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

package transport

import (
	"context"
	"errors"

	"github.com/nanobus/nanobus/channel"
	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/spec"
)

var ErrBadInput = errors.New("input was malformed")

type (
	NamedLoader func() (string, Loader)
	Loader      func(address string, namespaces spec.Namespaces, invoker Invoker, errorResolver errorz.Resolver, codecs ...channel.Codec) (Transport, error)

	Transport interface {
		Listen() error
		Close() error
	}

	Invoker func(ctx context.Context, namespace, service, id, function string, input interface{}) (interface{}, error)

	Registry map[string]Loader
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
