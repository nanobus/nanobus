package dapr

import (
	"context"
	"encoding/json"

	"github.com/dapr/dapr/pkg/actors"
	v1 "github.com/dapr/dapr/pkg/messaging/v1"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type InvokeActorConfig struct {
	// Key is the expression to evaluate the key to save.
	Type *expr.ValueExpr `mapstructure:"type"`
	// ForEach is an option expression to evaluate a
	ID *expr.ValueExpr `mapstructure:"id"`
	// Method is the name of the actor method to invoke.
	Method string `mapstructure:"method"`
	// Data is the optional data expression to tranform the data to set.
	Data *expr.DataExpr `mapstructure:"value"`
	// Metadata is the optional data expression for the key's metadata.
	Metadata *expr.DataExpr `mapstructure:"metadata"`
}

// InvokeActor is the NamedLoader for the Dapr get state operation
func InvokeActor() (string, actions.Loader) {
	return "@dapr/invoke_actor", InvokeActorLoader
}

func InvokeActorLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c InvokeActorConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var dapr *DaprComponents
	if err := resolve.Resolve(resolver,
		"dapr:components", &dapr); err != nil {
		return nil, err
	}

	return InvokeActorAction(dapr.Actors, &c), nil
}

func InvokeActorAction(
	actors actors.Actors,
	config *InvokeActorConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var payload interface{}
		req := v1.NewInvokeMethodRequest(config.Method)

		var err error
		if config.Data != nil {
			if payload, err = config.Data.Eval(data); err != nil {
				return nil, err
			}
		}
		if config.Metadata != nil {
			metadata, err := config.Metadata.EvalMap(data)
			if err != nil {
				return nil, err
			}
			md := make(map[string][]string, len(metadata))
			for k, v := range metadata {
				md[k] = []string{v}
			}

			req.WithMetadata(md)
		}

		if payload != nil {
			// TODO: encoding
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return nil, err
			}
			req.WithRawData(payloadBytes, "content/json")
		}

		response, err := actors.Call(ctx, req)
		if err != nil {
			return nil, err
		}

		var resp interface{}
		message := response.Message()
		if message != nil {
			if message.Data != nil && message.Data.Value != nil {
				// TODO: Decoding
				if err = json.Unmarshal(message.Data.Value, &resp); err != nil {
					return nil, err
				}
			}
		}

		return resp, nil
	}
}
