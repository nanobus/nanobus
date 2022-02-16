package dapr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type GetStateConfig struct {
	// Name is name of binding to invoke.
	Store string `mapstructure:"store" validate:"required"`
	// Operation is the name of the operation type for the binding to invoke
	Key *expr.ValueExpr `mapstructure:"key" validate:"required"`
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

	var httpClient HTTPClient
	if err := resolve.Resolve(resolver,
		"client:http", &httpClient); err != nil {
		return nil, err
	}

	return GetStateAction(httpClient, &c), nil
}

func GetStateAction(
	httpClient HTTPClient,
	config *GetStateConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		keyInt, err := config.Key.Eval(data)
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%v", keyInt)

		u, err := url.Parse(daprBaseURI)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "v1.0/state", config.Store, key)

		var response interface{}
		err = GET(ctx, httpClient,
			u.String(),
			func(data []byte) error {
				return json.Unmarshal(data, &response)
			})

		return response, err
	}
}
