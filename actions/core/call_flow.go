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

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
)

type CallFlowConfig struct {
	// Name is the name of the flow to call.
	Name string `mapstructure:"name" validate:"required"`
}

// Route is the NamedLoader for the filter action.
func CallFlow() (string, actions.Loader) {
	return "call_flow", CallFlowLoader
}

func CallFlowLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c CallFlowConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var processor Processor
	if err := resolve.Resolve(resolver,
		"system:processor", &processor); err != nil {
		return nil, err
	}

	return CallFlowAction(&c, processor), nil
}

func CallFlowAction(
	config *CallFlowConfig, processor Processor) actions.Action {
	return func(ctx context.Context, data actions.Data) (output interface{}, err error) {
		return processor.Flow(ctx, config.Name, data)
	}
}
