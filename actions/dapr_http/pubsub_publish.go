package dapr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type PublishMessageConfig struct {
	// Pubsub is name of pubsub to publish to.
	Pubsub string `mapstructure:"pubsub" validate:"required"`
	// Topic is the name of the topic to publish to.
	Topic string `mapstructure:"topic" validate:"required"`
	// Format is the format to write the data in.
	Format string `mapstructure:"format"`
	// Data is the input bindings sent
	Data *expr.DataExpr `mapstructure:"data"`
	// Metadata is the input binding metadata
	Metadata *expr.DataExpr `mapstructure:"metadata"`
}

// PublishMessage is the NamedLoader for Dapr pubsub publish message.
func PublishMessage() (string, actions.Loader) {
	return "@dapr/publish_message", PublishMessageLoader
}

func PublishMessageLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c PublishMessageConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var httpClient HTTPClient
	if err := resolve.Resolve(resolver,
		"client:http", &httpClient); err != nil {
		return nil, err
	}

	return PublishMessageAction(httpClient, &c), nil
}

func PublishMessageAction(
	httpClient HTTPClient,
	config *PublishMessageConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var err error

		var input interface{} = data
		if config.Data != nil {
			input, err = config.Data.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		var contentType string
		var requestBytes []byte

		switch config.Format {
		case "cloudevents+json":
			contentType = "application/cloudevents+json"
			if requestBytes, err = json.Marshal(input); err != nil {
				return nil, err
			}
		case "json":
			contentType = "application/json"
			if requestBytes, err = json.Marshal(input); err != nil {
				return nil, err
			}
		default:
			contentType = "application/json"
			if requestBytes, err = json.Marshal(input); err != nil {
				return nil, err
			}
		}

		u, err := url.Parse(daprBaseURI)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "v1.0/publish", config.Pubsub, config.Topic)

		req, err := http.NewRequestWithContext(
			ctx,
			"POST",
			u.String(),
			bytes.NewReader(requestBytes))
		req.Header.Set("Content-Type", contentType)
		if err != nil {
			return nil, err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			return nil, fmt.Errorf("expected 2XX status code; received %d", resp.StatusCode)
		}

		responseBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if len(responseBytes) > 0 {
			var response interface{}
			if err = json.Unmarshal(responseBytes, &response); err != nil {
				return nil, err
			}

			return response, nil
		}

		return nil, nil
	}
}
