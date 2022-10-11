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

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type FilterConfig struct {
	// Condition is the predicate expression for filtering.
	Condition *expr.ValueExpr `mapstructure:"condition" validate:"required"`
}

// Filter is the NamedLoader for the filter action.
func Filter() (string, actions.Loader) {
	return "filter", FilterLoader
}

func FilterLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c FilterConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return FilterAction(&c), nil
}

func FilterAction(
	config *FilterConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		resultInt, err := config.Condition.Eval(data)
		if err != nil {
			return nil, err
		}

		result, ok := resultInt.(bool)
		if !ok {
			return nil, fmt.Errorf("expression %q did not evaluate a boolean", config.Condition.Expr())
		}

		if !result {
			return nil, actions.Stop()
		}

		return nil, nil
	}
}
