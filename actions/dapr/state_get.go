package dapr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/components-contrib/state"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type GetStateConfig struct {
	// Name is name of binding to invoke.
	Store string `mapstructure:"store"`
	// Operation is the name of the operation type for the binding to invoke
	Key *expr.ValueExpr `mapstructure:"key"`
	// NotFoundError is the error to return if the key is not found
	NotFoundError string `mapstructure:"notFoundError"`
}

// GetState is the NamedLoader for the Dapr get state operation
func GetState() (string, actions.Loader) {
	return "@dapr/get_state", GetStateLoader
}

func GetStateLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c GetStateConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var dapr *DaprComponents
	if err := resolve.Resolve(resolver,
		"dapr:components", &dapr); err != nil {
		return nil, err
	}

	store, ok := dapr.StateStores[c.Store]
	if !ok {
		return nil, fmt.Errorf("state store %q not found", c.Store)
	}

	return GetStateAction(store, &c), nil
}

func GetStateAction(
	store state.Store,
	config *GetStateConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		keyInt, err := config.Key.Eval(data)
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%v", keyInt)

		resp, err := store.Get(&state.GetRequest{
			Key: key,
		})
		if err != nil {
			return nil, err
		}

		var response interface{}
		if len(resp.Data) > 0 {
			err = json.Unmarshal(resp.Data, &response)
		} else if config.NotFoundError != "" {
			return nil, errorz.Return(config.NotFoundError, errorz.Metadata{
				"store": config.Store,
				"key":   key,
			})
		}

		return response, err
	}
}
