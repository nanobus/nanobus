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
