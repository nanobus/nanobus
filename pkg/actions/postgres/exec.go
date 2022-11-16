package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/expr"
	"github.com/nanobus/nanobus/pkg/resolve"
	"github.com/nanobus/nanobus/pkg/resource"
)

type ExecConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `mapstructure:"resource" validate:"required"`
	// Data is the input bindings sent
	Data *expr.DataExpr `mapstructure:"data"`
	// SQL is the SQL query to execute.
	SQL string `mapstructure:"sql" validate:"required"`
	// Args are the evaluations to use as arguments for the SQL query.
	Args []*expr.ValueExpr `mapstructure:"args"`
}

// Exec is the NamedLoader for the invoke action.
func Exec() (string, actions.Loader) {
	return "@postgres/exec", ExecLoader
}

func ExecLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := ExecConfig{}
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

	return ExecAction(&c, pool), nil
}

func ExecAction(
	config *ExecConfig,
	pool *pgxpool.Pool) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var err error
		var input interface{} = map[string]interface{}(data)
		if config.Data != nil {
			input, err = config.Data.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		if multi, ok := input.([]interface{}); ok {
			if err = pool.BeginFunc(ctx, func(tx pgx.Tx) error {
				for _, item := range multi {
					if single, ok := item.(map[string]interface{}); ok {
						args := make([]interface{}, len(config.Args))
						for i, expr := range config.Args {
							var err error
							if args[i], err = expr.Eval(single); err != nil {
								return err
							}
						}

						_, err := tx.Exec(ctx, config.SQL, args...)
						if err != nil {
							return err
						}
						// if tag.RowsAffected() == 0 {
						// 	return errors.New("no rows effected")
						// }
					}
				}

				return nil
			}); err != nil {
				return nil, err
			}
		} else if single, ok := input.(map[string]interface{}); ok {
			args := make([]interface{}, len(config.Args))
			for i, expr := range config.Args {
				var err error
				if args[i], err = expr.Eval(single); err != nil {
					return nil, err
				}
			}

			_, err := pool.Exec(ctx, config.SQL, args...)
			if err != nil {
				return nil, err
			}
			// if tag.RowsAffected() == 0 {
			// 	return nil, errors.New("no rows effected")
			// }
		}

		return nil, nil
	}
}
