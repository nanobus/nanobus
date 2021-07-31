package actions

import (
	"context"
	"errors"

	"github.com/nanobus/nanobus/resolve"
)

type (
	Data        map[string]interface{}
	NamedLoader func() (string, Loader)
	Loader      func(with interface{}, resolver resolve.ResolveAs) (Action, error)
	Action      func(ctx context.Context, data Data) (interface{}, error)
	Registry    map[string]Loader
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}

func (d Data) Clone() Data {
	clone := make(Data, len(d)+5)
	for k, v := range d {
		clone[k] = v
	}
	return clone
}

// ErrStop is returned by an action when the processing should stop.
var ErrStop = errors.New("processing stopped")

func Stop() error {
	return ErrStop
}
