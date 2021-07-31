package core

import (
	"context"
	"log"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type LogConfig struct {
	Format string `mapstructure:"format"`
	// Args are the evaluations to use as arguments into the string format.
	Args []*expr.ValueExpr `mapstructure:"args"`
}

// Log is the NamedLoader for the log action.
func Log() (string, actions.Loader) {
	return "log", LogLoader
}

func LogLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c LogConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return LogAction(&c), nil
}

func LogAction(
	config *LogConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		args := make([]interface{}, len(config.Args))
		for i, expr := range config.Args {
			var err error
			if args[i], err = expr.Eval(data); err != nil {
				return nil, err
			}
		}

		log.Printf(config.Format, args...)

		return nil, nil
	}
}
