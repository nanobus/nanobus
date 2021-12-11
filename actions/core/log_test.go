package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/actions/core"
	"github.com/nanobus/nanobus/resolve"
)

type mockLogger struct {
	logr.Logger
	msg           string
	keysAndValues []interface{}
}

func (m *mockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.msg = msg
	m.keysAndValues = keysAndValues
}

func TestLog(t *testing.T) {
	ctx := context.Background()
	name, loader := core.Log()
	assert.Equal(t, "log", name)

	tests := []struct {
		name string

		config   map[string]interface{}
		resolver resolve.ResolveAs

		data      actions.Data
		format    string
		args      []interface{}
		loaderErr string
		actionErr string
	}{
		{
			name: "normal input",
			config: map[string]interface{}{
				"format": "testing %s %d %f",
				"args": []interface{}{
					"input.one",
					"input.two",
					"input.three",
				},
			},
			data: actions.Data{
				"input": map[string]interface{}{
					"one":   "1",
					"two":   2,
					"three": 3.0,
				},
			},
			format: "testing %s %d %f",
			args: []interface{}{
				"1",
				2,
				3.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logger mockLogger
			resolver := func(name string, target interface{}) bool {
				switch name {
				case "system:logger":
					return resolve.As(&logger, target)
				}
				return false
			}

			action, err := loader(tt.config, resolver)
			if tt.loaderErr != "" {
				require.EqualError(t, err, tt.loaderErr, "loader error was expected")
				return
			}
			require.NoError(t, err, "loader failed")

			_, err = action(ctx, tt.data)
			if tt.actionErr != "" {
				require.EqualError(t, err, tt.actionErr, "action error was expected")
				return
			}
			require.NoError(t, err, "action failed")
			assert.Equal(t, fmt.Sprintf(tt.format, tt.args...), logger.msg)
		})
	}
}
