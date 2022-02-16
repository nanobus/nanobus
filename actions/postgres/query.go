package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/resource"
	"github.com/nanobus/nanobus/stream"
)

type QueryConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `mapstructure:"resource" validate:"required"`
	// SQL is the SQL query to execute.
	SQL string `mapstructure:"sql" validate:"required"`
	// Args are the evaluations to use as arguments for the SQL query.
	Args []*expr.ValueExpr `mapstructure:"args"`
}

// Query is the NamedLoader for the invoke action.
func Query() (string, actions.Loader) {
	return "@postgres/query", QueryLoader
}

func QueryLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := QueryConfig{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var resources resource.Resources
	if err := resolve.Resolve(resolver,
		"resource:lookup", &resources); err != nil {
		return nil, err
	}

	poolI, ok := resources[c.Resource]
	if !ok {
		return nil, fmt.Errorf("resource %q is not registered", c.Resource)
	}
	pool, ok := poolI.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("resource %q is not a *pgxpool.Pool", c.Resource)
	}

	return QueryAction(&c, pool), nil
}

func QueryAction(
	config *QueryConfig,
	pool *pgxpool.Pool) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		s, ok := stream.FromContext(ctx)
		if !ok {
			return nil, errors.New("stream not in context")
		}

		args := make([]interface{}, len(config.Args))
		for i, expr := range config.Args {
			var err error
			if args[i], err = expr.Eval(data); err != nil {
				return nil, err
			}
		}

		rows, err := pool.Query(ctx, config.SQL, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		fields := rows.FieldDescriptions()
		fieldNames := make([]string, len(fields))
		for i, f := range fields {
			fieldNames[i] = snakeCaseToCamelCase(string(f.Name))
		}
		for rows.Next() {
			record := make(map[string]interface{})
			values, err := rows.Values()
			if err != nil {
				return nil, err
			}
			for i, v := range values {
				record[fieldNames[i]] = v
			}

			if err = s.SendData(record); err != nil {
				return nil, err
			}
		}

		return nil, nil
	}
}

func snakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false
	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToLower(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}
