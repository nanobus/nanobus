package dapr

import (
	"context"
	"fmt"

	"github.com/dapr/components-contrib/state"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type DeleteStateConfig struct {
	// Name is name of binding to invoke.
	Store string `mapstructure:"store" validate:"required"`
	// Operation is the name of the operation type for the binding to invoke.
	Key *expr.ValueExpr `mapstructure:"key" validate:"required"`
}

// DeleteState is the NamedLoader for the Dapr get state operation
func DeleteState() (string, actions.Loader) {
	return "@dapr/delete_state", DeleteStateLoader
}

func DeleteStateLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := DeleteStateConfig{}
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

	return DeleteStateAction(store, &c), nil
}

func DeleteStateAction(
	store state.Store,
	config *DeleteStateConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		keyInt, err := config.Key.Eval(data)
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%v", keyInt)

		err = store.Delete(&state.DeleteRequest{
			Key: key,
		})
		if err != nil {
			return nil, err
		}

		return nil, err
	}
}
