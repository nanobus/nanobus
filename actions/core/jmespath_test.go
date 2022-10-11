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
	"fmt"
	"testing"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/actions/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJMESPath(t *testing.T) {
	ctx := context.Background()
	name, loader := core.JMESPath()
	assert.Equal(t, "jmespath", name)

	tests := []struct {
		name      string
		config    map[string]interface{}
		data      actions.Data
		output    interface{}
		loaderErr string
		actionErr string
	}{
		{
			name: "normal input",
			config: map[string]interface{}{
				"path": `input.name`,
				"var":  `test`,
			},
			data: actions.Data{
				"input": map[string]interface{}{
					"name":        "test",
					"description": "full description",
					"nested": map[string]interface{}{
						"int":   1,
						"float": 1.1,
					},
				},
			},
			output: "test",
		},
		{
			name: "data input",
			config: map[string]interface{}{
				"path": `nested.int`,
				"data": `input`,
				"var":  `test`,
			},
			data: actions.Data{
				"input": map[string]interface{}{
					"name":        "test",
					"description": "full description",
					"nested": map[string]interface{}{
						"int":   1,
						"float": 1.1,
					},
				},
			},
			output: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := loader(ctx, tt.config, nil)
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
			assert.Equal(t, tt.output, output)
			if varName, ok := tt.config["var"]; ok {
				assert.Equal(t, tt.output, tt.data[fmt.Sprintf("%v", varName)])
			}
		})
	}
}
