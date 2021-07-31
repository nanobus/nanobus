package dapr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/components-contrib/pubsub"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type PublishMessageConfig struct {
	// Pubsub is name of pubsub to publish to.
	Pubsub string `mapstructure:"pubsub"`
	// Topic is the name of the topic to publish to.
	Topic string `mapstructure:"topic"`
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

	var dapr *DaprComponents
	if err := resolve.Resolve(resolver,
		"dapr:components", &dapr); err != nil {
		return nil, err
	}

	pubsub, ok := dapr.PubSubs[c.Pubsub]
	if !ok {
		return nil, fmt.Errorf("pubsub %q not found", c.Pubsub)
	}

	return PublishMessageAction(pubsub, &c), nil
}

func PublishMessageAction(
	component pubsub.PubSub,
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

		var requestBytes []byte
		// TODO: handle format
		if requestBytes, err = json.Marshal(input); err != nil {
			return nil, err
		}

		err = component.Publish(&pubsub.PublishRequest{
			Data:       requestBytes,
			PubsubName: config.Pubsub,
			Topic:      config.Topic,
		})

		return nil, err
	}
}
