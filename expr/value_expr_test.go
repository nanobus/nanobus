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

package expr_test

import (
	"testing"

	"github.com/nanobus/nanobus/expr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueExpr(t *testing.T) {
	var ve expr.ValueExpr
	err := ve.DecodeString(`input.test != nil && result.test == 5678`)
	require.NoError(t, err)
	result, err := ve.Eval(map[string]interface{}{
		"input": map[string]interface{}{
			"test": 1234,
		},
		"result": map[string]interface{}{
			"test": 5678,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, true, result)
}

var result interface{}

func BenchmarkEval(b *testing.B) {
	var ve expr.ValueExpr
	err := ve.DecodeString(`input.test != nil && result.test == 5678`)
	require.NoError(b, err)
	data := map[string]interface{}{
		"input": map[string]interface{}{
			"test": 1234,
		},
		"result": map[string]interface{}{
			"test": 5678,
		},
	}
	var r interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r, _ = ve.Eval(data)
	}
	result = r
}
