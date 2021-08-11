package core

import (
	"context"
	"fmt"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/codec"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
)

type DecodeConfig struct {
	TypeField string `mapstructure:"typeField"`
	DataField string `mapstructure:"dataField"`
	// Codec is the name of the codec to use for decoing.
	Codec string `mapstructure:"codec"`
	// Args are the arguments to pass to the decode function.
	CodecArgs []interface{} `mapstructure:"codecArgs"`
}

// Decode is the NamedLoader for the filter action.
func Decode() (string, actions.Loader) {
	return "decode", DecodeLoader
}

func DecodeLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := DecodeConfig{
		TypeField: "type",
		DataField: "input",
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
		current, ok := data[config.DataField]
		if !ok {
			return nil, nil
		}

		dataBytes, ok := current.([]byte)
		if !ok {
			return nil, fmt.Errorf("%q is not []byte which are required for decoding", config.DataField)
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

		return nil, nil
	}
}
