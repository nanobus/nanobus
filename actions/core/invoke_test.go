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
	"errors"
	"testing"

	"github.com/nanobus/iota/go/wasmrs/payload"
	"github.com/nanobus/iota/go/wasmrs/rx/mono"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/actions/core"
	"github.com/nanobus/nanobus/function"
	"github.com/nanobus/nanobus/resolve"
)

type mockInvoker struct {
	namespace string
	operation string
	input     payload.Payload
	output    payload.Payload
	err       error
}

func (m *mockInvoker) RequestResponse(ctx context.Context, namespace, operation string, p payload.Payload) mono.Mono[payload.Payload] {
	m.namespace = namespace
	m.operation = operation
	m.input = p
	if m.output != nil {
		return mono.Just(m.output)
	} else if p != nil {
		return mono.Just(p)
	}
	return mono.Error[payload.Payload](m.err)
}

func TestInvoke(t *testing.T) {
	ctx := context.Background()
	name, loader := core.Invoke()
	assert.Equal(t, "invoke", name)

	tests := []struct {
		name string

		invoker  *mockInvoker
		config   map[string]interface{}
		resolver resolve.ResolveAs

		data      actions.Data
		namespace string
		operation string
		output    interface{}
		loaderErr string
		actionErr string
	}{
		{
			name:    "normal input",
			invoker: &mockInvoker{},
			config: map[string]interface{}{
				"namespace": "test.v1",
				"operation": "test",
			},
			data: actions.Data{
				"input": map[string]interface{}{
					"test": "test",
				},
			},
			output: map[string]interface{}{
				"test": "test",
			},
			namespace: "test.v1",
			operation: "test",
		},
		{
			name:    "normal input",
			invoker: &mockInvoker{},
			config: map[string]interface{}{
				"input": `{
					"test": input.test + "12345",
				}`,
				"namespace": "test.v1",
				"operation": "test",
			},
			data: actions.Data{
				"input": map[string]interface{}{
					"test": "test",
				},
			},
			output: map[string]interface{}{
				"test": "test12345",
			},
			namespace: "test.v1",
			operation: "test",
		},
		{
			name:    "bytes input",
			invoker: &mockInvoker{},
			config: map[string]interface{}{
				"namespace": "test.v1",
				"operation": "test",
			},
			data: actions.Data{
				"input": []byte(`{ "test": "test" }`),
			},
			output: map[string]interface{}{
				"test": "test",
			},
			namespace: "test.v1",
			operation: "test",
		},
		{
			name:    "string input",
			invoker: &mockInvoker{},
			config: map[string]interface{}{
				"namespace": "test.v1",
				"operation": "test",
			},
			data: actions.Data{
				"input": `{ "test": "test" }`,
			},
			output: map[string]interface{}{
				"test": "test",
			},
			namespace: "test.v1",
			operation: "test",
		},
		{
			name:      "invoke from context",
			invoker:   &mockInvoker{},
			config:    map[string]interface{}{},
			data:      actions.Data{},
			namespace: "test.v1",
			operation: "test",
		},
		{
			name: "invoke error",
			invoker: &mockInvoker{
				err: errors.New("test error"),
			},
			config:    map[string]interface{}{},
			data:      actions.Data{},
			namespace: "test.v1",
			operation: "test",
			actionErr: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := func(name string, target interface{}) bool {
				switch name {
				case "client:invoker":
					return resolve.As(tt.invoker, target)
				}
				return false
			}

			action, err := loader(ctx, tt.config, resolver)
			if tt.loaderErr != "" {
				require.EqualError(t, err, tt.loaderErr, "loader error was expected")
				return
			}
			require.NoError(t, err, "loader failed")

			fctx := function.ToContext(ctx, function.Function{
				Namespace: tt.namespace,
				Operation: tt.operation,
			})
			output, err := action(fctx, tt.data)
			if tt.actionErr != "" {
				require.EqualError(t, err, tt.actionErr, "action error was expected")
				return
			}
			require.NoError(t, err, "action failed")
			assert.Equal(t, tt.namespace, tt.invoker.namespace)
			assert.Equal(t, tt.operation, tt.invoker.operation)
			assert.Equal(t, tt.output, output)
		})
	}
}
