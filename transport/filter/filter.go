package filter

import (
	"context"

	"github.com/nanobus/nanobus/resolve"
)

type (
	NamedLoader func() (string, Loader)
	Loader      func(with interface{}, resolver resolve.ResolveAs) (Filter, error)
	Filter      func(ctx context.Context, header Header) (context.Context, error)
	Registry    map[string]Loader

	Header interface {
		Get(name string) string
		Values(name string) []string
	}
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
