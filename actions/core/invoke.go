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
	"encoding/json"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/channel"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/function"
	"github.com/nanobus/nanobus/resolve"
)

type InvokeConfig struct {
	// Namespace of the service to invoke.
	Namespace string `mapstructure:"namespace"`
	// Operation of the service to invoke.
	Operation string `mapstructure:"operation"`
	// Input optionally transforms the input sent to the function.
	Input *expr.DataExpr `mapstructure:"input"`
}

type Invoker interface {
	InvokeWithReturn(ctx context.Context, receiver channel.Receiver, input, output interface{}) error
}

// Invoke is the NamedLoader for the invoke action.
func Invoke() (string, actions.Loader) {
	return "invoke", InvokeLoader
}

func InvokeLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := InvokeConfig{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var invoker Invoker
	if err := resolve.Resolve(resolver,
		"client:invoker", &invoker); err != nil {
		return nil, err
	}

	return InvokeAction(invoker, &c), nil
}

func InvokeAction(
	invoker Invoker,
	config *InvokeConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		input := data["input"]
		if config.Input != nil {
			var err error
			input, err = config.Input.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		switch v := input.(type) {
		case []byte:
			if err := json.Unmarshal(v, &input); err != nil {
				return nil, err
			}
		case string:
			if err := json.Unmarshal([]byte(v), &input); err != nil {
				return nil, err
			}
		}

		namespace := config.Namespace
		operation := config.Operation
		namespaceEmpty := namespace == ""
		operationEmpty := operation == ""

		// Grab the incoming function details if needed.
		if namespaceEmpty || operationEmpty {
			fn := function.FromContext(ctx)

			if namespaceEmpty {
				namespace = fn.Namespace
			}
			if operationEmpty {
				operation = fn.Operation
			}
		}

		var response interface{}
		if err := invoker.InvokeWithReturn(ctx, channel.Receiver{
			Namespace: namespace,
			Operation: operation,
		}, input, &response); err != nil {
			return nil, err
		}
		if response != nil {
			return response, nil
		}

		return nil, nil
	}
}
