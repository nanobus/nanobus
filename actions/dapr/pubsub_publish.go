package dapr

import (
	"context"
	"fmt"

	"github.com/dapr/components-contrib/pubsub"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/codec"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

type PublishMessageConfig struct {
	// Pubsub is name of pubsub to publish to.
	Pubsub string `mapstructure:"pubsub"`
	// Topic is the name of the topic to publish to.
	Topic string `mapstructure:"topic"`
	// Codec is the configured codec to use for encoding the message.
	Codec string `mapstructure:"codec"`
	// CodecArgs are the arguments for the codec, if any.
	CodecArgs []interface{} `mapstructure:"codecArgs"`
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
	var codecs codec.Codecs
	if err := resolve.Resolve(resolver,
		"dapr:components", &dapr,
		"codec:lookup", &codecs); err != nil {
		return nil, err
	}

	codec, ok := codecs[c.Codec]
	if !ok {
		return nil, fmt.Errorf("codec %q not found", c.Codec)
	}

	pubsub, ok := dapr.PubSubs[c.Pubsub]
	if !ok {
		return nil, fmt.Errorf("pubsub %q not found", c.Pubsub)
	}

	return PublishMessageAction(pubsub, codec, &c), nil
}

func PublishMessageAction(
	component pubsub.PubSub,
	codec codec.Codec,
	config *PublishMessageConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var err error

		var input interface{} = data["input"]
		if config.Data != nil {
			input, err = config.Data.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		dataBytes, err := codec.Encode(input, config.CodecArgs...)
		if err != nil {
			return nil, err
		}

		err = component.Publish(&pubsub.PublishRequest{
			Data:       dataBytes,
			PubsubName: config.Pubsub,
			Topic:      config.Topic,
			Metadata: map[string]string{
				"rawPayload": "true",
			},
		})

		return nil, err
	}
}
