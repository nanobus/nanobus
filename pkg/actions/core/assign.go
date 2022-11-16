package core

import (
	"context"

	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/coalesce"
	"github.com/nanobus/nanobus/pkg/expr"
	"github.com/nanobus/nanobus/pkg/resolve"
)

type AssignConfig struct {
	Value *expr.ValueExpr `mapstructure:"value"`
	Data  *expr.DataExpr  `mapstructure:"data"`
	To    string          `mapstructure:"to"`
}

// Assign is the NamedLoader for the assign action.
func Assign() (string, actions.Loader) {
	return "assign", AssignLoader
}

func AssignLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c AssignConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return AssignAction(&c), nil
}

func AssignAction(
	config *AssignConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (output interface{}, err error) {
		if config.Value != nil {
			output, err = config.Value.Eval(data)
			if err != nil {
				return nil, err
			}
		} else if config.Data != nil {
			output, err = config.Data.Eval(data)
			if err != nil {
				return nil, err
			}
			if v, ok := coalesce.ToMapSI(output, true); ok {
				output = v
			}
		}

		if config.To != "" {
			data[config.To] = output
		}

		return output, nil
	}
}
