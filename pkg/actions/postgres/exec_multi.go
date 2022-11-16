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

type ExecMultiConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `mapstructure:"resource" validate:"required"`
	// Statements are the statements to execute within a single transaction.
	Statements []Statement `mapstructure:"statements"`
}

type Statement struct {
	// Data is the input bindings sent
	Data *expr.DataExpr `mapstructure:"data"`
	// SQL is the SQL query to execute.
	SQL string `mapstructure:"sql" validate:"required"`
	// Args are the evaluations to use as arguments for the SQL query.
	Args []*expr.ValueExpr `mapstructure:"args"`
}

// ExecMulti is the NamedLoader for the invoke action.
func ExecMulti() (string, actions.Loader) {
	return "@postgres/exec_multi", ExecMultiLoader
}

func ExecMultiLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := ExecMultiConfig{}
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

	return ExecMultiAction(&c, pool), nil
}

func ExecMultiAction(
	config *ExecMultiConfig,
	pool *pgxpool.Pool) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		err := pool.BeginFunc(ctx, func(tx pgx.Tx) error {
			for _, stmt := range config.Statements {
				var err error
				var input interface{} = map[string]interface{}(data)
				if stmt.Data != nil {
					input, err = stmt.Data.Eval(data)
					if err != nil {
						return err
					}
				}

				if multi, ok := input.([]interface{}); ok {
					for _, item := range multi {
						if single, ok := item.(map[string]interface{}); ok {
							single["$root"] = data
							args := make([]interface{}, len(stmt.Args))
							for i, expr := range stmt.Args {
								var err error
								if args[i], err = expr.Eval(single); err != nil {
									delete(single, "$root")
									return err
								}
							}

							_, err := tx.Exec(ctx, stmt.SQL, args...)
							if err != nil {
								delete(single, "$root")
								return err
							}
							// if tag.RowsAffected() == 0 {
							// 	delete(single, "$root")
							// 	return errors.New("no rows effected")
							// }
							delete(single, "$root")
						}
					}

					return nil
				} else if single, ok := input.(map[string]interface{}); ok {
					single["$root"] = data
					args := make([]interface{}, len(stmt.Args))
					for i, expr := range stmt.Args {
						var err error
						if args[i], err = expr.Eval(single); err != nil {
							delete(single, "$root")
							return err
						}
					}

					_, err := tx.Exec(ctx, stmt.SQL, args...)
					if err != nil {
						delete(single, "$root")
						return err
					}
					// if tag.RowsAffected() == 0 {
					// 	delete(single, "$root")
					// 	return errors.New("no rows effected")
					// }
					delete(single, "$root")
				}
			}
			return nil
		})

		return nil, err
	}
}
