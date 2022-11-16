package wasmrs

import (
	"context"
	"os"

	"github.com/nanobus/iota/go/wasmrs/host"

	"github.com/nanobus/nanobus/pkg/compute"
	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
	"github.com/nanobus/nanobus/pkg/runtime"
)

type Config struct {
	// Filename is the file name of the WasmRS module to load.
	// TODO: Load from external location
	Filename runtime.FilePath `mapstructure:"filename" validate:"required"`
}

// WasmRS
func WasmRS() (string, compute.Loader) {
	return "wasmrs", Loader
}

func Loader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (compute.Invoker, error) {
	c := Config{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	source, err := os.ReadFile(string(c.Filename))
	if err != nil {
		return nil, err
	}

	h, err := host.New(ctx)
	if err != nil {
		return nil, err
	}
	module, err := h.Compile(ctx, source)
	if err != nil {
		return nil, err
	}

	return module.Instantiate(ctx)
}
