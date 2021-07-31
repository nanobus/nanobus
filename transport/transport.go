package transport

import (
	"context"
	"errors"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/nanobus/spec"
)

var ErrBadInput = errors.New("input was malformed")

type (
	NamedLoader func() (string, Loader)
	Loader      func(address string, namespaces spec.Namespaces, invoker Invoker, codecs ...functions.Codec) (Transport, error)

	Transport interface {
		Listen() error
		Close() error
	}

	Invoker func(ctx context.Context, namespace, service, function string, input interface{}) (interface{}, error)

	Registry map[string]Loader
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
