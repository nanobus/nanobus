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

package registry

import (
	"context"

	"github.com/nanobus/nanobus/resolve"
)

type (
	NamedLoader[T any] func() (string, Loader[T])
	Loader[T any]      func(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (T, error)
	Registry[T any]    map[string]Loader[T]
)

func (r Registry[T]) Register(loaders ...NamedLoader[T]) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}
