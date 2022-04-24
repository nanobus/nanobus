/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type LogConfig struct {
	Format string `mapstructure:"format" validate:"required"`
	// Args are the evaluations to use as arguments into the string format.
	Args []*expr.ValueExpr `mapstructure:"args"`
}

type Logger interface {
	Printf(format string, v ...interface{})
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

	var logger logr.Logger
	if err := resolve.Resolve(resolver,
		"system:logger", &logger); err != nil {
		return nil, err
	}

	return LogAction(logger, &c), nil
}

func LogAction(
	logger logr.Logger,
	config *LogConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		args := make([]interface{}, len(config.Args))
		for i, expr := range config.Args {
			var err error
			if args[i], err = expr.Eval(data); err != nil {
				return nil, err
			}
		}

		msg := config.Format
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}

		logger.Info(msg)

		return nil, nil
	}
}
