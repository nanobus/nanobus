package filter

import (
	"context"
	"net/http"

	"github.com/nanobus/nanobus/resolve"
)

type (
	NamedLoader func() (string, Loader)
	Loader      func(with interface{}, resolver resolve.ResolveAs) (Filter, error)
	Filter      func(ctx context.Context, req *http.Request) (context.Context, error)
	Registry    map[string]Loader
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
