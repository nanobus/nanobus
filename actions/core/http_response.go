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

	"github.com/cenkalti/backoff/v4"
	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/transport/httpresponse"
)

type HTTPResponseConfig struct {
	Status  int                  `mapstructure:"status"`
	Headers []HTTPResponseHeader `mapstructure:"headers"`
}

type HTTPResponseHeader struct {
	Name  string          `mapstructure:"name" validate:"required"`
	Value *expr.ValueExpr `mapstructure:"value" validate:"required"`
}

// HTTPResponse is the NamedLoader for Dapr output bindings
func HTTPResponse() (string, actions.Loader) {
	return "http_response", HTTPResponseLoader
}

func HTTPResponseLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := HTTPResponseConfig{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return HTTPResponseAction(&c), nil
}

func HTTPResponseAction(
	config *HTTPResponseConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		resp := httpresponse.FromContext(ctx)
		if resp == nil {
			return nil, nil
		}

		if config.Status > 0 {
			resp.Status = config.Status
		}

		for _, h := range config.Headers {
			val, err := h.Value.Eval(data)
			if err != nil {
				return nil, backoff.Permanent(err)
			}
			resp.Header.Add(h.Name, fmt.Sprintf("%v", val))
		}

		return nil, nil
	}
}
