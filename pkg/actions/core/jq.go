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

	"github.com/itchyny/gojq"

	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
)

func JQLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c JQConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	query, err := gojq.Parse(c.Query)
	if err != nil {
		return nil, fmt.Errorf("error parsing jq query: %w", err)
	}

	code, err := gojq.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("error compiling jq query: %w", err)
	}

	return JQAction(&c, code), nil
}

func JQAction(
	config *JQConfig, code *gojq.Code) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var in interface{} = map[string]interface{}(data)
		if config.Data != nil {
			var err error
			in, err = config.Data.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		var emitted []interface{}
		iter := code.RunWithContext(ctx, in)
		for {
			out, ok := iter.Next()
			if !ok {
				break
			}

			if err, ok := out.(error); ok {
				return nil, err
			}

			if config.Single {
				if config.Var != nil {
					data[*config.Var] = out
				}
				return out, nil
			}

			emitted = append(emitted, out)
		}

		if config.Single {
			return nil, nil
		}

		if config.Var != nil {
			data[*config.Var] = emitted
		}

		return emitted, nil
	}
}
