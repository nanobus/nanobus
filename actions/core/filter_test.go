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

package core_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/actions/core"
	"github.com/nanobus/nanobus/resolve"
)

func TestFilter(t *testing.T) {
	ctx := context.Background()
	name, loader := core.Filter()
	assert.Equal(t, "filter", name)

	tests := []struct {
		name string

		config   map[string]interface{}
		resolver resolve.ResolveAs

		data      actions.Data
		expected  interface{}
		loaderErr string
		actionErr string
	}{
		{
			name: "continue",
			config: map[string]interface{}{
				"condition": "test == true",
			},
			data: actions.Data{
				"test": true,
			},
		},
		{
			name: "stop",
			config: map[string]interface{}{
				"condition": "test == false",
			},
			data: actions.Data{
				"test": true,
			},
			actionErr: actions.ErrStop.Error(),
		},
		{
			name: "non-boolean expression",
			config: map[string]interface{}{
				"condition": "12345",
			},
			data: actions.Data{
				"test": true,
			},
			actionErr: "expression \"12345\" did not evaluate a boolean",
		},
		{
			name: "loader error",
			config: map[string]interface{}{
				"condition": 12345,
			},
			loaderErr: "1 error(s) decoding:\n\n* 'condition' expected a map, got 'int'",
		},
		{
			name: "expression error",
			config: map[string]interface{}{
				"condition": "test.test == true",
			},
			data:      actions.Data{},
			actionErr: "cannot fetch test from <nil> (1:6)\n | test.test == true\n | .....^",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := loader(tt.config, tt.resolver)
			if tt.loaderErr != "" {
				require.EqualError(t, err, tt.loaderErr, "loader error was expected")
				return
			}

			require.NoError(t, err, "loader failed")

			output, err := action(ctx, tt.data)
			if tt.actionErr != "" {
				require.EqualError(t, err, tt.actionErr, "action error was expected")
				return
			}
			require.NoError(t, err, "action failed")
			assert.Equal(t, tt.expected, output)
		})
	}
}
