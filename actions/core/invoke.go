package core

import (
	"context"
	"encoding/json"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type InvokeConfig struct {
	// Namespace of the service to invoke.
	Namespace string `mapstructure:"namespace"`
	// Operation of the service to invoke.
	Operation string `mapstructure:"operation"`
	// Input optionally transforms the input sent to the function.
	Input *expr.DataExpr `mapstructure:"input"`
}

// Invoke is the NamedLoader for the invoke action.
func Invoke() (string, actions.Loader) {
	return "invoke", InvokeLoader
}

func InvokeLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := InvokeConfig{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var invoker *functions.Invoker
	if err := resolve.Resolve(resolver,
		"client:invoker", &invoker); err != nil {
		return nil, err
	}

	return InvokeAction(invoker, &c), nil
}

func InvokeAction(
	invoker *functions.Invoker,
	config *InvokeConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		input := data["input"]
		if config.Input != nil {
			var err error
			input, err = config.Input.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		if inputBytes, ok := input.([]byte); ok {
			if err := json.Unmarshal(inputBytes, &input); err != nil {
				return nil, err
			}
		}
		if inputString, ok := input.(string); ok {
			if err := json.Unmarshal([]byte(inputString), &input); err != nil {
				return nil, err
			}
		}

		var response interface{}
		if err := invoker.InvokeWithReturn(ctx, config.Namespace, config.Operation, input, &response); err != nil {
			return nil, err
		}
		if response != nil {
			return response, nil
		}

		return nil, nil
	}
}
