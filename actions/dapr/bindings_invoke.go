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
	"fmt"

	"github.com/dapr/components-contrib/bindings"

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

	var dapr *DaprComponents
	if err := resolve.Resolve(resolver,
		"dapr:components", &dapr); err != nil {
		return nil, err
	}

	binding, ok := dapr.OutputBindings[c.Name]
	if !ok {
		return nil, fmt.Errorf("output binding %q not found", c.Name)
	}

	return InvokeBindingAction(binding, &c), nil
}

func InvokeBindingAction(
	binding bindings.OutputBinding,
	config *InvokeBindingConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var bindingData interface{}
		r := bindings.InvokeRequest{
			Operation: bindings.OperationKind(config.Operation),
		}

		var err error
		if config.Data != nil {
			if bindingData, err = config.Data.Eval(data); err != nil {
				return nil, err
			}
		}
		if config.Metadata != nil {
			if r.Metadata, err = config.Metadata.EvalMap(data); err != nil {
				return nil, err
			}
		}

		// TODO: multi-codec support
		if r.Data, err = json.Marshal(bindingData); err != nil {
			return nil, err
		}

		resp, err := binding.Invoke(&r)
		if err != nil {
			return nil, err
		}

		var response interface{}
		if len(resp.Data) > 0 {
			err = coalesce.JSONUnmarshal(resp.Data, &response)
		}

		return response, err
	}
}
