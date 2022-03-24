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

type SetStateConfig struct {
	// Store is name of state store to invoke.
	Store string `mapstructure:"store" validate:"required"`
	// Key is the expression to evaluate the key to save.
	Key *expr.ValueExpr `mapstructure:"key" validate:"required"`
	// ForEach is an option expression to evaluate a
	ForEach *expr.ValueExpr `mapstructure:"forEach"`
	// Value is the optional data expression to tranform the data to set.
	Value *expr.DataExpr `mapstructure:"value"`
	// Metadata is the optional data expression for the key's metadata.
	Metadata *expr.DataExpr `mapstructure:"metadata"`
}

// SetState is the NamedLoader for the Dapr get state operation
func SetState() (string, actions.Loader) {
	return "@dapr/set_state", SetStateLoader
}

func SetStateLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c SetStateConfig
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

	return SetStateAction(store, &c), nil
}

type SetItem struct {
	Key      string            `json:"key"`
	Value    interface{}       `json:"value"`
	Etag     *string           `json:"etag,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func SetStateAction(
	store state.Store,
	config *SetStateConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var items []interface{}
		if config.ForEach != nil {
			itemsInt, err := config.ForEach.Eval(data)
			if err != nil {
				return nil, fmt.Errorf("could not evaluate data: %w", err)
			}
			var ok bool
			if items, ok = itemsInt.([]interface{}); ok {
				return nil, fmt.Errorf("forEach expression %q did not return a slice of items", config.ForEach.Expr())
			}
		}

		var r []state.SetRequest
		if items == nil {
			it, err := createSetItem(data, nil, config)
			if err != nil {
				return nil, err
			}

			r = []state.SetRequest{it}
		} else {
			r = make([]state.SetRequest, len(items))
			for i, item := range items {
				it, err := createSetItem(data, item, config)
				if err != nil {
					return nil, err
				}

				r[i] = it
			}
		}

		err := store.BulkSet(r)

		return nil, err
	}
}

func createSetItem(
	data actions.Data,
	item interface{},
	config *SetStateConfig) (it state.SetRequest, err error) {
	variables := make(map[string]interface{}, len(data)+1)
	for k, v := range data {
		variables[k] = v
	}
	variables["item"] = item

	it = state.SetRequest{
		Value: variables["input"],
	}
	keyInt, err := config.Key.Eval(variables)
	if err != nil {
		return it, fmt.Errorf("could not evaluate key: %w", err)
	}
	it.Key = fmt.Sprintf("%v", keyInt)

	if config.Value != nil {
		if it.Value, err = config.Value.Eval(variables); err != nil {
			return it, fmt.Errorf("could not evaluate value: %w", err)
		}
	}
	if config.Metadata != nil {
		if it.Metadata, err = config.Metadata.EvalMap(variables); err != nil {
			return it, fmt.Errorf("could not evaluate metadata: %w", err)
		}
	}

	return it, nil
}
