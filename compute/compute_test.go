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

package compute_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nanobus/nanobus/compute"
	"github.com/nanobus/nanobus/resolve"
)

func TestRegistry(t *testing.T) {
	r := compute.Registry{}

	loader := func(with interface{}, resolver resolve.ResolveAs) (compute.Invoker, error) {
		return nil, nil
	}
	namedLoader := func() (string, compute.Loader) {
		return "test", loader
	}

	r.Register(namedLoader)

	assert.Equal(t, fmt.Sprintf("%v", compute.Loader(loader)), fmt.Sprintf("%p", r["test"]))
}
