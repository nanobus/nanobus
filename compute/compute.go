package compute

import (
	"context"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/nanobus/resolve"
)

type (
	BusInvoker   func(ctx context.Context, namespace, service, function string, input interface{}) (interface{}, error)
	StateInvoker func(ctx context.Context, namespace, id, key string) ([]byte, error)
	NamedLoader  func() (string, Loader)
	Loader       func(with interface{}, resolver resolve.ResolveAs) (*Compute, error)
	Registry     map[string]Loader

	Compute struct {
		Invoker           *functions.Invoker
		Start             func() error
		WaitUntilShutdown func() error
		Close             func() error
		Environ           func() []string
	}
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
