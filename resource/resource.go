package resource

import (
	"context"

	"github.com/nanobus/nanobus/resolve"
)

type (
	NamedLoader func() (string, Loader)
	Loader      func(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (interface{}, error)
	Registry    map[string]Loader
	Resources   map[string]interface{}
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
