package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/runtime"
)

type SetStateConfig struct {
	// Store is name of state store to invoke.
	Store string `mapstructure:"store" validate:"required"`
	// Key is the expression to evaluate the key to save.
	Key *expr.ValueExpr `mapstructure:"key" validate:"required"`
	// ForEach is an option expression to evaluate a
	ForEach *expr.ValueExpr `mapstructure:"forEach"`
	// Value is the optional data expression to tranform the data to set.
	Value *expr.DataExpr `mapstructure:"value" validate:"required"`
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

	var httpClient HTTPClient
	var env runtime.Environment
	if err := resolve.Resolve(resolver,
		"client:http", &httpClient,
		"os:env", &env); err != nil {
		return nil, err
	}

	return SetStateAction(httpClient, &c), nil
}

type SetItem struct {
	Key      string            `json:"key"`
	Value    interface{}       `json:"value"`
	Etag     *string           `json:"etag,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func SetStateAction(
	httpClient HTTPClient,
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

		var r []SetItem
		if items == nil {
			it, err := createSetItem(data, nil, config)
			if err != nil {
				return nil, err
			}

			r = []SetItem{it}
		} else {
			r = make([]SetItem, len(items))
			for i, item := range items {
				it, err := createSetItem(data, item, config)
				if err != nil {
					return nil, err
				}

				r[i] = it
			}
		}

		u, err := url.Parse(daprBaseURI)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "v1.0/state", config.Store)

		var response interface{}
		err = POST(ctx, httpClient,
			u.String(),
			func() ([]byte, error) {
				return json.Marshal(&r)
			}, func(data []byte) error {
				return coalesce.JSONUnmarshal(data, &response)
			})

		return response, err
	}
}

func createSetItem(
	data actions.Data,
	item interface{},
	config *SetStateConfig) (it SetItem, err error) {
	variables := make(map[string]interface{}, len(data)+1)
	for k, v := range data {
		variables[k] = v
	}
	variables["item"] = item

	it = SetItem{
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
