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

package dapr

import (
	"context"
	"encoding/json"
	"net/url"
	"path"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type InvokeBindingConfig struct {
	// Name is name of binding to invoke.
	Name string `mapstructure:"name" validate:"required"`
	// Operation is the name of the operation type for the binding to invoke
	Operation string `mapstructure:"operation" validate:"required"`
	// Data is the input bindings sent
	Data *expr.DataExpr `mapstructure:"data"`
	// Metadata is the input binding metadata
	Metadata *expr.DataExpr `mapstructure:"metadata"`
}

// InvokeBinding is the NamedLoader for Dapr output bindings
func InvokeBinding() (string, actions.Loader) {
	return "@dapr/invoke_binding", InvokeBindingLoader
}

func InvokeBindingLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c InvokeBindingConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var httpClient HTTPClient
	if err := resolve.Resolve(resolver,
		"client:http", &httpClient); err != nil {
		return nil, err
	}

	return InvokeBindingAction(httpClient, &c), nil
}

func InvokeBindingAction(
	httpClient HTTPClient,
	config *InvokeBindingConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		type requestBody struct {
			Operation string            `json:"operation"`
			Data      interface{}       `json:"data,omitempty"`
			Metadata  map[string]string `json:"metadata,omitempty"`
		}

		r := requestBody{
			Operation: config.Operation,
			Data:      data,
		}

		var err error
		if config.Data != nil {
			if r.Data, err = config.Data.Eval(data); err != nil {
				return nil, err
			}
		}
		if config.Metadata != nil {
			if r.Metadata, err = config.Metadata.EvalMap(data); err != nil {
				return nil, err
			}
		}

		u, err := url.Parse(daprBaseURI)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "v1.0/bindings", config.Name)

		var response interface{}
		err = POST(ctx, httpClient,
			u.String(),
			func() ([]byte, error) {
				return json.Marshal(&r)
			}, func(data []byte) error {
				return coalesce.JSONUnmarshal(data, &response)
			})

		return response, err
	}
}
