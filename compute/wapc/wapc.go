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

package wapc

import (
	"context"
	"errors"
	"fmt"
	"os"
	go_runtime "runtime"
	"strings"

	"github.com/nanobus/go-functions"
	wapc_mux "github.com/nanobus/go-functions/transports/wapc"
	wapc "github.com/wapc/wapc-go"
	"github.com/wapc/wapc-go/engines/wazero"

	"github.com/nanobus/nanobus/compute"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/resolve"
)

type WaPCConfig struct {
	// Filename is the file name of the waPC/WebAssembly module to load.
	Filename string `mapstructure:"filename"` // TODO: Load from external location
	// PoolSize controls the number of waPC instance of the module to create and pool.
	// It also represents the maximum number of concurrent requests the module can process.
	PoolSize uint64 `mapstructure:"poolSize"`
}

// Mux is the NamedLoader for the waPC compute.
func WaPC() (string, compute.Loader) {
	return "wapc", WaPCLoader
}

func WaPCLoader(with interface{}, resolver resolve.ResolveAs) (*compute.Compute, error) {
	var busInvoker compute.BusInvoker
	var msgpackcodec functions.Codec
	if err := resolve.Resolve(resolver,
		"bus:invoker", &busInvoker,
		"codec:msgpack", &msgpackcodec); err != nil {
		return nil, err
	}

	c := WaPCConfig{
		PoolSize: uint64(go_runtime.NumCPU() * 5),
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	wasmBytes, err := os.ReadFile(c.Filename)
	if err != nil {
		return nil, err
	}

	engine := wazero.Engine()

	module, err := engine.New(wasmBytes, func(ctx context.Context, binding, namespace, operation string, payload []byte) ([]byte, error) {
		lastDot := strings.LastIndexByte(namespace, '.')
		if lastDot < 0 {
			return nil, fmt.Errorf("invalid namespace %q", namespace)
		}
		service := namespace[lastDot+1:]
		namespace = namespace[:lastDot]

		var input interface{}
		if err := msgpackcodec.Decode(payload, &input); err != nil {
			return nil, err
		}

		result, err := busInvoker(ctx, namespace, service, operation, input)
		if err != nil {
			return nil, err
		}

		return msgpackcodec.Encode(result)
	})
	if err != nil {
		return nil, err
	}

	module.SetLogger(wapc.Println)
	module.SetWriter(wapc.Print)

	m, err := wapc_mux.New(module, uint64(c.PoolSize))
	if err != nil {
		return nil, err
	}
	invoke := func(ctx context.Context, receiver functions.Receiver, payload []byte) ([]byte, error) {
		resp, err := m.Invoke(ctx, receiver, payload)
		if err != nil {
			// Trim out wrapped message.
			msg := err.Error()
			msg = strings.TrimPrefix(msg, "Host error: ")
			i := strings.Index(msg, "; ~lib/@wapc/")
			if i > 0 {
				msg = msg[:i]
			}
			return nil, errors.New(msg)
		}
		return resp, nil
	}
	invokeStream := func(ctx context.Context, receiver functions.Receiver) (functions.Streamer, error) {
		return nil, errorz.New(errorz.Unimplemented, "streaming is not implemented for waPC")
	}
	invoker := functions.NewInvoker(invoke, invokeStream, msgpackcodec)
	done := make(chan struct{}, 1)

	return &compute.Compute{
		Invoker: invoker,
		Start:   func() error { return nil },
		WaitUntilShutdown: func() error {
			<-done
			return nil
		},
		Close: func() error {
			close(done)
			return nil
		},
	}, nil
}
