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

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/stream"
)

type TestConfig struct {
	Data *expr.DataExpr `mapstructure:"data"`
}

// Test is the NamedLoader for the log action.
func Test() (string, actions.Loader) {
	return "@postgres/test", TestLoader
}

func TestLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c TestConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return TestAction(&c), nil
}

func TestAction(
	config *TestConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		s, ok := stream.SinkFromContext(ctx)
		if !ok {
			return nil, errors.New("stream not in context")
		}

		v, err := config.Data.Eval(data)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for i := 0; i < 10; i++ {
			if err = s.Next(v, nil); err != nil {
				fmt.Println(err)
				return nil, err
			}
		}

		return nil, nil
	}
}
