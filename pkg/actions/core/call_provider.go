/*
 * Copyright 2022 The NanoBus Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package core

import (
	"context"
	"fmt"

	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
)

func CallProviderLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c CallProviderConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var processor Processor
	if err := resolve.Resolve(resolver,
		"system:processor", &processor); err != nil {
		return nil, err
	}

	return CallProviderAction(&c, processor), nil
}

func CallProviderAction(
	config *CallProviderConfig, processor Processor) actions.Action {
	return func(ctx context.Context, data actions.Data) (output interface{}, err error) {
		if config.Input != nil {
			input, err := config.Input.Eval(data)
			if err != nil {
				return nil, err
			}
			data = data.Clone()
			data["input"] = input
			data["pipe"] = input
			data["$"] = input
		}

		var ok bool
		output, ok, err = processor.Provider(ctx, config.Handler, data)
		if !ok {
			return nil, fmt.Errorf("could not find provider %s", &config.Handler)
		}
		return
	}
}
