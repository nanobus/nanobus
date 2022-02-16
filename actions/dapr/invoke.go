package dapr

import (
	"context"
	"encoding/json"

	"github.com/dapr/dapr/pkg/messaging"
	invokev1 "github.com/dapr/dapr/pkg/messaging/v1"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type InvokeConfig struct {
	// Name is name of SQL binding to invoke.
	AppID  string `mapstructure:"appId" validate:"required"`
	Method string `mapstructure:"method" validate:"required"`
	Verb   string `mapstructure:"verb" validate:"required"`
	// Data is the input bindings sent
	Data *expr.DataExpr `mapstructure:"data"`
}

// Invoke is the NamedLoader for Dapr output bindings
func Inoke() (string, actions.Loader) {
	return "@dapr/invoke", InokeLoader
}

func InokeLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c InvokeConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var dapr *DaprComponents
	if err := resolve.Resolve(resolver,
		"dapr:components", &dapr); err != nil {
		return nil, err
	}

	return InokeAction(dapr.DirectMessaging, &c), nil
}

func InokeAction(
	directMessaging messaging.DirectMessaging,
	config *InvokeConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var err error
		body := interface{}(data)
		if config.Data != nil {
			if body, err = config.Data.Eval(data); err != nil {
				return nil, err
			}
		}

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		// Construct internal invoke method request
		req := invokev1.NewInvokeMethodRequest(config.Method).WithHTTPExtension(config.Verb, "")
		req.WithRawData(bodyBytes, "application/json")
		req.WithMetadata(map[string][]string{})
		// Save headers to internal metadata
		//req.WithFastHTTPHeaders(&reqCtx.Request.Header)

		resp, err := directMessaging.Invoke(ctx, config.AppID, req)
		if err != nil {
			return nil, err
		}

		_, rawBody := resp.RawData()

		var response interface{}
		if len(rawBody) > 0 {
			err = json.Unmarshal(rawBody, &response)
		}

		return response, err
	}
}
