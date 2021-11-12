package codec_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nanobus/nanobus/codec"
	"github.com/nanobus/nanobus/resolve"
)

func TestRegistry(t *testing.T) {
	r := codec.Registry{}

	loader := func(with interface{}, resolver resolve.ResolveAs) (codec.Codec, error) {
		return nil, nil
	}
	namedLoader := func() (string, codec.Loader) {
		return "test", loader
	}

	r.Register(namedLoader)

	assert.Equal(t, fmt.Sprintf("%v", codec.Loader(loader)), fmt.Sprintf("%p", r["test"]))
}
