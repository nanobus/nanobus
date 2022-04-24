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

package function_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nanobus/nanobus/function"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	empty := function.FromContext(ctx)
	assert.Equal(t, function.Function{}, empty)
	fn := function.Function{
		Namespace: "test.v1",
		Operation: "testing",
	}
	fctx := function.ToContext(ctx, fn)
	actual := function.FromContext(fctx)
	assert.Equal(t, fn, actual)
}
