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

package structerror_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nanobus/nanobus/structerror"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		create   func() *structerror.Error
		metadata map[string]string
		str      string
	}{
		{
			name: "new",
			create: func() *structerror.Error {
				return structerror.New("not_found",
					"key", "abcdef",
					"store", "statestore")
			},
			metadata: map[string]string{
				"key":   "abcdef",
				"store": "statestore",
			},
			str: `not_found
[key] abcdef
[store] statestore`,
		},
		{
			name: "parse",
			create: func() *structerror.Error {
				contents := `not_found
ignore
[key] abcdef
[store] statestore`
				return structerror.Parse(contents)
			},
			metadata: map[string]string{
				"key":   "abcdef",
				"store": "statestore",
			},
			str: `not_found
[key] abcdef
[store] statestore`,
		},
		{
			name: "parse no metadata",
			create: func() *structerror.Error {
				contents := `not_found`
				return structerror.Parse(contents)
			},
			metadata: nil,
			str:      `not_found`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.create()
			assert.Equal(t, "not_found", e.Code())
			assert.Equal(t, tt.metadata, e.Metadata())

			assert.Equal(t, tt.str, e.Error())
			assert.Equal(t, tt.str, e.String())
		})
	}
}
