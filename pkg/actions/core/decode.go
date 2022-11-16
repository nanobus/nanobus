package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/codec"
	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
)

type DecodeConfig struct {
	TypeField string `mapstructure:"typeField"`
	DataField string `mapstructure:"dataField"`
	// Codec is the name of the codec to use for decoing.
	Codec string `mapstructure:"codec" validate:"required"`
	// Args are the arguments to pass to the decode function.
	CodecArgs []interface{} `mapstructure:"codecArgs"`
}

// Decode is the NamedLoader for the filter action.
func Decode() (string, actions.Loader) {
	return "decode", DecodeLoader
}

func DecodeLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := DecodeConfig{
		TypeField: "type",
		DataField: "input",
		Codec:     "json",
		CodecArgs: []interface{}{},
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var codecs codec.Codecs
	if err := resolve.Resolve(resolver,
		"codec:lookup", &codecs); err != nil {
		return nil, err
	}

	codec, ok := codecs[c.Codec]
	if !ok {
		return nil, fmt.Errorf("unknown codec %q", c.Codec)
	}

	return DecodeAction(codec, &c), nil
}

func DecodeAction(
	codec codec.Codec,
	config *DecodeConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		parts := strings.Split(config.DataField, ".")
		var current interface{} = map[string]interface{}(data)
		for _, part := range parts {
			var ok bool
			asMap, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("non-map encountered for property %q", part)
			}
			current, ok = asMap[part]
			if !ok {
				return nil, fmt.Errorf("property %q not set", part)
			}
		}

		var dataBytes []byte
		switch v := current.(type) {
		case []byte:
			dataBytes = v
		case string:
			dataBytes = []byte(v)
		default:
			return nil, fmt.Errorf("%q must be []byte or string for decoding", config.DataField)
		}

		decoded, typeName, err := codec.Decode(dataBytes, config.CodecArgs...)
		if err != nil {
			return nil, err
		}

		if typeName != "" && config.TypeField != "" {
			data[config.TypeField] = typeName
		}
		if config.DataField != "" {
			data[config.DataField] = decoded
		}

		return decoded, nil
	}
}
