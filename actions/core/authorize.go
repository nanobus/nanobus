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
	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/security/claims"
)

type AuthorizeConfig struct {
	// Condition is the predicate expression for authorization.
	Condition *expr.ValueExpr        `mapstructure:"condition"`
	Has       []string               `mapstructure:"has"`
	Check     map[string]interface{} `mapstructure:"check"`
	Error     string                 `mapstructure:"error"`
}

// Authorize is the NamedLoader for the log action.
func Authorize() (string, actions.Loader) {
	return "authorize", AuthorizeLoader
}

func AuthorizeLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := AuthorizeConfig{
		Error: "permission_denied",
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return AuthorizeAction(&c), nil
}

func AuthorizeAction(
	config *AuthorizeConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		if config.Condition != nil {
			resultInt, err := config.Condition.Eval(data)
			if err != nil {
				return nil, err
			}

			result, ok := resultInt.(bool)
			if !ok {
				return nil, fmt.Errorf("expression %q did not evaluate a boolean", config.Condition.Expr())
			}

			if !result {
				return nil, errorz.Return(config.Error, errorz.Metadata{
					"expr": config.Condition.Expr(),
				})
			}
		}

		claimsMap := claims.FromContext(ctx)

		for _, claim := range config.Has {
			if _, ok := claimsMap[claim]; !ok {
				return nil, errorz.Return(config.Error, errorz.Metadata{
					"claim": claim,
				})
			}
		}

		for claim, value := range config.Check {
			v := claimsMap[claim]
			if v != value {
				return nil, errorz.Return(config.Error, errorz.Metadata{
					"claim": claim,
					"want":  value,
				})
			}
		}

		return nil, nil
	}
}
