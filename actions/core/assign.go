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

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
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

func AssignLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
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
