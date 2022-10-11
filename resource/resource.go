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

package resource

import (
	"fmt"
	"reflect"

	"github.com/nanobus/nanobus/registry"
)

type (
	NamedLoader = registry.NamedLoader[any]
	Loader      = registry.Loader[any]
	Registry    = registry.Registry[any]

	// NamedLoader func() (string, Loader)
	// Loader      func(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (interface{}, error)
	// Registry    map[string]Loader
	Resources map[string]interface{}
)

// func (r Registry) Register(loaders ...NamedLoader) {
// 	for _, l := range loaders {
// 		name, loader := l()
// 		r[name] = loader
// 	}
// }

func Get[T any](r Resources, name string) (res T, err error) {
	var iface interface{}
	iface, ok := r[name]
	if !ok {
		return res, fmt.Errorf("resource %q is not registered", name)
	}
	res, ok = iface.(T)
	if !ok {
		return res, fmt.Errorf("resource %q is not a %s", name, reflect.TypeOf(res).Name())
	}

	return res, nil
}
