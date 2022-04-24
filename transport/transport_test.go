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

package transport_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/spec"
	"github.com/nanobus/nanobus/transport"
)

func TestRegistry(t *testing.T) {
	r := transport.Registry{}

	loader := func(address string, namespaces spec.Namespaces, invoker transport.Invoker, errorResolver errorz.Resolver, codecs ...functions.Codec) (transport.Transport, error) {
		return nil, nil
	}
	namedLoader := func() (string, transport.Loader) {
		return "test", loader
	}

	r.Register(namedLoader)

	assert.Equal(t, fmt.Sprintf("%v", transport.Loader(loader)), fmt.Sprintf("%p", r["test"]))
}
