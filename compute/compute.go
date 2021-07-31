package compute

import (
	"context"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/nanobus/resolve"
)

type (
	BusInvoker  func(ctx context.Context, namespace, service, function string, input interface{}) (interface{}, error)
	NamedLoader func() (string, Loader)
	Loader      func(with interface{}, resolver resolve.ResolveAs) (*functions.Invoker, error)
	Registry    map[string]Loader
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
