package core_test

import (
	"context"
	"errors"
	"testing"

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
	input     interface{}
	output    interface{}
	err       error
}

func (m *mockInvoker) InvokeWithReturn(ctx context.Context, namespace, operation string, input interface{}, output interface{}) error {
	m.namespace = namespace
	m.operation = operation
	m.input = input
	if m.output != nil {
		resolve.As(m.output, output)
	} else if input != nil {
		resolve.As(input, output)
	}
	return m.err
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

			action, err := loader(tt.config, resolver)
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