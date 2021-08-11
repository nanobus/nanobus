package codec

import (
	"github.com/nanobus/nanobus/resolve"
)

type (
	// Codec is an interface that handles encoding and decoding payloads send to and
	// received from functions.
	Codec interface {
		ContentType() string
		Encode(v interface{}, args ...interface{}) ([]byte, error)
		Decode(data []byte, args ...interface{}) (interface{}, string, error)
	}

	NamedLoader func() (string, Loader)
	Loader      func(with interface{}, resolver resolve.ResolveAs) (Codec, error)
	Registry    map[string]Loader
	Codecs      map[string]Codec
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
