/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package wasmrs

import (
	"context"
	"os"

	"github.com/WasmRS/wasmrs-go/host"
	"github.com/nanobus/nanobus/compute"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
)

type Config struct {
	// Filename is the file name of the WasmRS module to load.
	Filename string `mapstructure:"filename" validate:"required"` // TODO: Load from external location
}

// WasmRS
func WasmRS() (string, compute.Loader) {
	return "wasmrs", Loader
}

func Loader(with interface{}, resolver resolve.ResolveAs) (compute.Invoker, error) {
	ctx := context.Background()
	c := Config{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	source, err := os.ReadFile(c.Filename)
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
