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

package compute

import (
	"context"
	"io"

	"github.com/WasmRS/wasmrs-go/invoke"
	"github.com/WasmRS/wasmrs-go/operations"
	"github.com/WasmRS/wasmrs-go/payload"
	"github.com/WasmRS/wasmrs-go/rx/flux"
	"github.com/WasmRS/wasmrs-go/rx/mono"

	"github.com/nanobus/nanobus/resolve"
)

type (
	BusInvoker   func(ctx context.Context, namespace, service, function string, input interface{}) (interface{}, error)
	StateInvoker func(ctx context.Context, namespace, id, key string) ([]byte, error)
	NamedLoader  func() (string, Loader)
	Loader       func(with interface{}, resolver resolve.ResolveAs) (Invoker, error)
	Registry     map[string]Loader

	// Compute struct {
	// 	Invoker           Invoker
	// 	Start             func() error
	// 	WaitUntilShutdown func() error
	// 	Close             func() error
	// 	Environ           func() []string
	// }

	Invoker interface {
		io.Closer
		Operations() operations.Table

		FireAndForget(context.Context, payload.Payload)
		RequestResponse(context.Context, payload.Payload) mono.Mono[payload.Payload]
		RequestStream(context.Context, payload.Payload) flux.Flux[payload.Payload]
		RequestChannel(context.Context, payload.Payload, flux.Flux[payload.Payload]) flux.Flux[payload.Payload]

		SetRequestResponseHandler(index uint32, handler invoke.RequestResponseHandler)
		SetFireAndForgetHandler(index uint32, handler invoke.FireAndForgetHandler)
		SetRequestStreamHandler(index uint32, handler invoke.RequestStreamHandler)
		SetRequestChannelHandler(index uint32, handler invoke.RequestChannelHandler)
	}
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
