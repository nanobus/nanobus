package json

import (
	"encoding/json"

	"github.com/nanobus/nanobus/codec"
	"github.com/nanobus/nanobus/resolve"
)

type (
	// Codec encodes and decodes Avro records.
	Codec struct{}
)

// JSON is the NamedLoader for this codec.
func JSON() (string, codec.Loader) {
	return "json", Loader
}

func Loader(with interface{}, resolver resolve.ResolveAs) (codec.Codec, error) {
	return NewCodec(), nil
}

// NewCodec creates a `Codec`.
func NewCodec() *Codec {
	return &Codec{}
}

func (c *Codec) ContentType() string {
	return "application/json"
}

// Decode decodes JSON bytes to a value.
func (c *Codec) Decode(msgValue []byte, args ...interface{}) (interface{}, string, error) {
	var data interface{}
	if err := json.Unmarshal(msgValue, &data); err != nil {
		return nil, "", err
	}

	return data, "", nil
}

// Encode encodes a value into JSON encoded bytes.
func (c *Codec) Encode(value interface{}, args ...interface{}) ([]byte, error) {
	return json.Marshal(value)
}
