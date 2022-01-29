package json

import (
	"encoding/json"

	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/codec"
	"github.com/nanobus/nanobus/resolve"
)

type (
	// Codec encodes and decodes Avro records.
	Codec struct{}
)

// CloudEventsJSON is the NamedLoader for this codec.
func CloudEventsJSON() (string, bool, codec.Loader) {
	return "cloudevents+json", true, Loader
}

func Loader(with interface{}, resolver resolve.ResolveAs) (codec.Codec, error) {
	return NewCodec(), nil
}

// NewCodec creates a `Codec`.
func NewCodec() *Codec {
	return &Codec{}
}

func (c *Codec) ContentType() string {
	return "application/cloudevents+json"
}

// Decode decodes JSON bytes to a value.
func (c *Codec) Decode(msgValue []byte, args ...interface{}) (interface{}, string, error) {
	var data map[string]interface{}
	if err := coalesce.JSONUnmarshal(msgValue, &data); err != nil {
		return nil, "", err
	}

	var typeValue string
	if typeField, ok := data["type"]; ok {
		typeValue, _ = typeField.(string)
	}

	return data, typeValue, nil
}

// Encode encodes a value into JSON encoded bytes.
func (c *Codec) Encode(value interface{}, args ...interface{}) ([]byte, error) {
	return json.Marshal(value)
}
