package core

import (
	"context"
	"fmt"

	"github.com/itchyny/gojq"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type JQConfig struct {
	// Query is the predicate expression for filtering.
	Query string `mapstructure:"query" validate:"required"`
	// Data is the optional data expression to pass to jq.
	Data *expr.DataExpr `mapstructure:"data"`
	// Single, if true, returns the first result.
	Single bool `mapstructure:"single"`
	// Var, if set, is the variable that is set with the result.
	Var string `mapstructure:"var"`
}

// JQ is the NamedLoader for the jq action.
func JQ() (string, actions.Loader) {
	return "jq", JQLoader
}

func JQLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c JQConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	query, err := gojq.Parse(c.Query)
	if err != nil {
		return nil, fmt.Errorf("error parsing jq query: %w", err)
	}

	code, err := gojq.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("error compiling jq query: %w", err)
	}

	return JQAction(&c, code), nil
}

func JQAction(
	config *JQConfig, code *gojq.Code) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var in interface{} = map[string]interface{}(data)
		if config.Data != nil {
			var err error
			in, err = config.Data.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		var emitted []interface{}
		iter := code.RunWithContext(ctx, in)
		for {
			out, ok := iter.Next()
			if !ok {
				break
			}

			if err, ok := out.(error); ok {
				return nil, err
			}

			if config.Single {
				if config.Var != "" {
					data[config.Var] = out
				}
				return out, nil
			}

			emitted = append(emitted, out)
		}

		if config.Single {
			return nil, nil
		}

		if config.Var != "" {
			data[config.Var] = emitted
		}

		return emitted, nil
	}
}
