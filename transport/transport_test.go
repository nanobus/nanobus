package transport_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/nanobus/spec"
	"github.com/nanobus/nanobus/transport"
)

func TestRegistry(t *testing.T) {
	r := transport.Registry{}

	loader := func(address string, namespaces spec.Namespaces, invoker transport.Invoker, codecs ...functions.Codec) (transport.Transport, error) {
		return nil, nil
	}
	namedLoader := func() (string, transport.Loader) {
		return "test", loader
	}

	r.Register(namedLoader)

	assert.Equal(t, fmt.Sprintf("%v", transport.Loader(loader)), fmt.Sprintf("%p", r["test"]))
}
