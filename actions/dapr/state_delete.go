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
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/resource"
)

type DeleteStateConfig struct {
	// Name is name of binding to invoke.
	Resource string `mapstructure:"resource" validate:"required"`
	// Operation is the name of the operation type for the binding to invoke.
	Key *expr.ValueExpr `mapstructure:"key" validate:"required"`
}

// DeleteState is the NamedLoader for the Dapr get state operation
func DeleteState() (string, actions.Loader) {
	return "@dapr/delete_state", DeleteStateLoader
}

func DeleteStateLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := DeleteStateConfig{}
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

	return DeleteStateAction(client, &c), nil
}

func DeleteStateAction(
	client proto.StateStoreClient,
	config *DeleteStateConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		keyInt, err := config.Key.Eval(data)
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%v", keyInt)

		_, err = client.Delete(ctx, &proto.DeleteRequest{
			Key: key,
		})
		if err != nil {
			return nil, err
		}

		return nil, err
	}
}
