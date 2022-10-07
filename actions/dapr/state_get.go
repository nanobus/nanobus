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
	"fmt"

	proto "github.com/dapr/dapr/pkg/proto/components/v1"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/resource"
)

type GetStateConfig struct {
	// Resource is name of binding to invoke.
	Resource string `mapstructure:"resource" validate:"required"`
	// Operation is the name of the operation type for the binding to invoke.
	Key *expr.ValueExpr `mapstructure:"key" validate:"required"`
	// NotFoundError is the error to return if the key is not found.
	NotFoundError string `mapstructure:"notFoundError"`
	// Var, if set, is the variable that is set with the result.
	Var string `mapstructure:"var"`
}

// GetState is the NamedLoader for the Dapr get state operation
func GetState() (string, actions.Loader) {
	return "@dapr/get_state", GetStateLoader
}

func GetStateLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := GetStateConfig{
		NotFoundError: "not_found",
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var resources resource.Resources
	if err := resolve.Resolve(resolver,
		"resource:lookup", &resources); err != nil {
		return nil, err
	}

	client, err := resource.Get[proto.StateStoreClient](resources, c.Resource)
	if err != nil {
		return nil, err
	}

	return GetStateAction(client, &c), nil
}

func GetStateAction(
	client proto.StateStoreClient,
	config *GetStateConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		keyInt, err := config.Key.Eval(data)
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%v", keyInt)

		resp, err := client.Get(ctx, &proto.GetRequest{
			Key: key,
			// TODO
		})
		if err != nil {
			return nil, err
		}

		var response interface{}
		if len(resp.Data) > 0 {
			// TODO: use codec
			err = coalesce.JSONUnmarshal(resp.Data, &response)
		} else if config.NotFoundError != "" {
			return nil, errorz.Return(config.NotFoundError, errorz.Metadata{
				"resource": config.Resource,
				"key":      key,
			})
		}

		if config.Var != "" {
			data[config.Var] = response
		}

		return response, err
	}
}
