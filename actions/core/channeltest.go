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

package core

import (
	"context"
	"errors"
	"io"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/stream"
)

type ChannelTestConfig struct {
}

// ChannelTest is the NamedLoader for the channel test action.
func ChannelTest() (string, actions.Loader) {
	return "channel_test", ChannelTestLoader
}

func ChannelTestLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := ChannelTestConfig{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return ChannelTestAction(&c), nil
}

func ChannelTestAction(
	config *ChannelTestConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		source, ok := stream.SourceFromContext(ctx)
		if !ok {
			return nil, errors.New("stream not in context")
		}
		sink, ok := stream.SinkFromContext(ctx)
		if !ok {
			return nil, errors.New("stream not in context")
		}

		input, _ := data["input"].(map[string]interface{})
		prefix, _ := input["prefix"].(string)

		var in map[string]interface{}
		for {
			if err := source.Next(&in, nil); err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}

			out := make(map[string]interface{}, len(in))
			for k, v := range in {
				switch v := v.(type) {
				case string:
					out[k] = prefix + v
				default:
					out[k] = v
				}
			}

			if err := sink.Next(out, nil); err != nil {
				return nil, err
			}
		}

		return nil, nil
	}
}
