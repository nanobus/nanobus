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

package filter

import (
	"context"

	"github.com/nanobus/nanobus/registry"
)

type (
	NamedLoader = registry.NamedLoader[Filter]
	Loader      = registry.Loader[Filter]
	Registry    = registry.Registry[Filter]

	Filter func(ctx context.Context, header Header) (context.Context, error)

	Header interface {
		Get(name string) string
		Values(name string) []string
	}
)
