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
	"errors"
	"fmt"

	"github.com/dapr/dapr/pkg/actors"
	v1 "github.com/dapr/dapr/pkg/messaging/v1"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type InvokeActorConfig struct {
	// Type is the expression to evaluate the key to save.
	Type string `mapstructure:"type" validate:"required"`
	// ForEach is an option expression to evaluate a
	ID *expr.ValueExpr `mapstructure:"id" validate:"required"`
	// Method is the name of the actor method to invoke.
	Method string `mapstructure:"method" validate:"required"`
	// Data is the optional data expression to tranform the data to set.
	Data *expr.DataExpr `mapstructure:"data"`
	// Metadata is the optional data expression for the key's metadata.
	Metadata *expr.DataExpr `mapstructure:"metadata"`
}

// InvokeActor is the NamedLoader for the Dapr get state operation
func InvokeActor() (string, actions.Loader) {
	return "@dapr/invoke_actor", InvokeActorLoader
}

func InvokeActorLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c InvokeActorConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var dapr *DaprComponents
	if err := resolve.Resolve(resolver,
		"dapr:components", &dapr); err != nil {
		return nil, err
	}

	if dapr.Actors == nil {
		return nil, errors.New("actor system not initialized")
	}

	return InvokeActorAction(dapr.Actors, &c), nil
}

func InvokeActorAction(
	actors actors.Actors,
	config *InvokeActorConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		idInt, err := config.ID.Eval(data)
		if err != nil {
			return nil, err
		}
		id := fmt.Sprintf("%v", idInt)

		var payload interface{} = data["input"]
		req := v1.
			NewInvokeMethodRequest(config.Method).
			WithActor(config.Type, id)

		if config.Data != nil {
			if payload, err = config.Data.Eval(data); err != nil {
				return nil, err
			}
		}
		if config.Metadata != nil {
			metadata, err := config.Metadata.EvalMap(data)
			if err != nil {
				return nil, err
			}
			md := make(map[string][]string, len(metadata))
			for k, v := range metadata {
				md[k] = []string{v}
			}

			req.WithMetadata(md)
		}

		if payload != nil {
			// TODO: encoding
			payloadBytes, err := msgpack.Marshal(payload)
			if err != nil {
				return nil, err
			}
			req.WithRawData(payloadBytes, "application/msgpack")
		}

		response, err := actors.Call(ctx, req)
		if err != nil {
			return nil, err
		}

		var resp interface{}
		message := response.Message()
		if message != nil {
			if message.Data != nil && message.Data.Value != nil {
				// TODO: Decoding
				if err = json.Unmarshal(message.Data.Value, &resp); err != nil {
					return nil, err
				}
			}
		}

		return resp, nil
	}
}
