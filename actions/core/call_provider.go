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
	"errors"
	"strings"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
)

type CallProviderConfig struct {
	// Namespace is the namespace of the provider to call.
	Namespace string `mapstructure:"namespace" validate:"required"`
	// Operation is the operation name of the provider to call.
	Operation string `mapstructure:"operation" validate:"required"`
}

// Route is the NamedLoader for the filter action.
func CallProvider() (string, actions.Loader) {
	return "call_provider", CallProviderLoader
}

func CallProviderLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c CallProviderConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var processor Processor
	if err := resolve.Resolve(resolver,
		"system:processor", &processor); err != nil {
		return nil, err
	}

	namespace := c.Namespace
	i := strings.LastIndex(namespace, ".")
	if i < 0 {
		return nil, errors.New("invalid namespace")
	}
	service := namespace[i+1:]
	namespace = namespace[:i]

	return CallProviderAction(namespace, service, c.Operation, processor), nil
}

func CallProviderAction(
	namespace, service, operation string, processor Processor) actions.Action {
	return func(ctx context.Context, data actions.Data) (output interface{}, err error) {
		return processor.Provider(ctx, namespace, service, operation, data)
	}
}
