/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	Pubsub string `mapstructure:"pubsub" validate:"required"`
	// Topic is the name of the topic to publish to.
	Topic string `mapstructure:"topic" validate:"required"`
	// Codec is the configured codec to use for encoding the message.
	Codec string `mapstructure:"codec" validate:"required"`
	// CodecArgs are the arguments for the codec, if any.
	CodecArgs []interface{} `mapstructure:"codecArgs"`
	// Key is the optional value to use for the message key (is supported).
	Key *expr.ValueExpr `mapstructure:"key"`
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

		var key string
		if config.Key != nil {
			keyInt, err := config.Key.Eval(data)
			if err != nil {
				return nil, err
			}
			key = fmt.Sprintf("%v", keyInt)
		}

		dataBytes, err := codec.Encode(input, config.CodecArgs...)
		if err != nil {
			return nil, err
		}

		metadata := map[string]string{
			"rawPayload": "true",
		}

		if key != "" {
			metadata["partitionKey"] = key
		}

		err = component.Publish(&pubsub.PublishRequest{
			Data:       dataBytes,
			PubsubName: config.Pubsub,
			Topic:      config.Topic,
			Metadata:   metadata,
		})

		return nil, err
	}
}
