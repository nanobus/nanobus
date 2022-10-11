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
	"github.com/nanobus/nanobus/runtime"
)

type runner func(ctx context.Context, data actions.Data) (interface{}, error)

type mockProcessor struct {
	pipeline  *runtime.Pipeline
	runnables map[string]runner
	ran       []string
	err       error
}

var _ = (core.Processor)((*mockProcessor)(nil))

func (m *mockProcessor) LoadPipeline(pl *runtime.Pipeline) (runtime.Runnable, error) {
	m.pipeline = pl
	fn := m.runnables[pl.Name]

	runnable := mockRunnable{m, pl.Name, fn}

	return runnable.Run, m.err
}

func (m *mockProcessor) Pipeline(ctx context.Context, name string, data actions.Data) (interface{}, error) {
	return data, nil
}

func (m *mockProcessor) Provider(ctx context.Context, namespace, service, function string, data actions.Data) (interface{}, error) {
	return data, nil
}

func (m *mockProcessor) Event(ctx context.Context, name string, data actions.Data) (interface{}, error) {
	return data, nil
}

type mockRunnable struct {
	m       *mockProcessor
	summary string
	fn      runner
}

func (m mockRunnable) Run(ctx context.Context, data actions.Data) (interface{}, error) {
	m.m.ran = append(m.m.ran, m.summary)
	return m.fn(ctx, data)
}

func TestRoute(t *testing.T) {
	ctx := context.Background()
	name, loader := core.Route()
	assert.Equal(t, "route", name)

	tests := []struct {
		name string

		config    map[string]interface{}
		processor *mockProcessor

		data      actions.Data
		pipeline  runtime.Pipeline
		expected  interface{}
		ran       []string
		loaderErr string
		actionErr string
	}{
		{
			name: "single",
			config: map[string]interface{}{
				"selection": "single",
				"routes": []interface{}{
					map[string]interface{}{
						"name": "A",
						"when": `path == 'A'`,
						"then": []interface{}{
							map[string]interface{}{
								"name": "1",
								"uses": "test a",
							},
						},
					},
					map[string]interface{}{
						"name": "B",
						"when": `path == 'B'`,
						"then": []interface{}{
							map[string]interface{}{
								"name": "1",
								"uses": "test b",
							},
						},
					},
				},
			},
			processor: &mockProcessor{
				runnables: map[string]runner{
					"B": func(ctx context.Context, data actions.Data) (interface{}, error) {
						return "b", nil
					},
				},
			},
			pipeline: runtime.Pipeline{
				Name: "B",
				Steps: []runtime.Step{
					{
						Name: "1",
						Uses: "test b",
					},
				},
			},
			expected: "b",
			ran:      []string{"B"},
			data: actions.Data{
				"path": "B",
			},
		},
		{
			name: "multi",
			config: map[string]interface{}{
				"selection": "multi",
				"routes": []interface{}{
					map[string]interface{}{
						"name": "A",
						"when": `path == 'A'`,
						"then": []interface{}{
							map[string]interface{}{
								"name": "1",
								"uses": "test a",
							},
						},
					},
					map[string]interface{}{
						"name": "B",
						"when": `other == 'B'`,
						"then": []interface{}{
							map[string]interface{}{
								"name": "1",
								"uses": "test b",
							},
						},
					},
				},
			},
			processor: &mockProcessor{
				runnables: map[string]runner{
					"A": func(ctx context.Context, data actions.Data) (interface{}, error) {
						return "a", nil
					},
					"B": func(ctx context.Context, data actions.Data) (interface{}, error) {
						return "b", nil
					},
				},
			},
			pipeline: runtime.Pipeline{
				Name: "B",
				Steps: []runtime.Step{
					{
						Name: "1",
						Uses: "test b",
					},
				},
			},
			expected: "b",
			ran:      []string{"A", "B"},
			data: actions.Data{
				"path":  "A",
				"other": "B",
			},
		},
		{
			name: "configuration error",
			config: map[string]interface{}{
				"selection": "invalid",
				"routes":    1234,
			},
			data:      actions.Data{},
			loaderErr: "2 error(s) decoding:\n\n* 'routes': source data must be an array or slice, got int\n* error decoding 'selection': invalid SelectionMode \"invalid\": unexpected selection mode: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := func(name string, target interface{}) bool {
				switch name {
				case "system:processor":
					return resolve.As(tt.processor, target)
				}
				return false
			}

			action, err := loader(ctx, tt.config, resolver)
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
			assert.Equal(t, tt.ran, tt.processor.ran)
			assert.Equal(t, &tt.pipeline, tt.processor.pipeline)
			assert.Equal(t, tt.expected, output)
		})
	}
}
