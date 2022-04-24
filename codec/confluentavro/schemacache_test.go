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

package confluentavro_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nanobus/nanobus/codec/confluentavro"
)

type (
	mockSchemaRegistry struct {
		name    string
		counter int
	}
)

var (
	schemaString = `{
	"type": "record",
	"namespace": "com.example",
	"name": "FullName",
	"fields": [{
			"name": "first",
			"type": "string"
		},
		{
			"name": "last",
			"type": "string"
		}
	]
}`
	happySchema, _ = confluentavro.ParseSchema(0, schemaString)
)

func (m *mockSchemaRegistry) GetSchemaByID(id int) (string, error) {
	m.counter++
	switch m.name {
	case "schema-retrieval-error":
		return "", errors.New("woops")
	case "happy-path":
		return schemaString, nil
	case "schema-parse-error":
		return "bad", nil
	default:
		return "", nil
	}
}

func TestSchemaCache_GetSchema(t *testing.T) {
	// normal test cases
	subtests := []struct {
		name   string
		schema *confluentavro.Schema
		err    error
	}{
		{
			name:   "schema-retrieval-error",
			err:    errors.New("error getting schema ID 0: woops"),
			schema: nil,
		},
		{
			name:   "schema-parse-error",
			err:    errors.New("error parsing schema ID 0: avro: unknown type: bad"),
			schema: nil,
		},
		{
			name:   "happy-path",
			err:    nil,
			schema: happySchema,
		},
	}

	for _, tt := range subtests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockSchemaRegistry{name: tt.name}
			schemaCache := confluentavro.NewSchemaCache(&mock, 10)
			schema, err := schemaCache.GetSchema(0)
			if err != nil {
				if tt.err.Error() != err.Error() {
					t.Errorf("expected error (%v), got error (%v)", tt.err, err)
				}
			}
			if !reflect.DeepEqual(schema, tt.schema) {
				t.Errorf("expected schema (%v), got schema (%v)", tt.schema, schema)
			}
		})
	}

	// cache hit test cases
	mock := mockSchemaRegistry{name: "happy-path"}
	schemaCache := confluentavro.NewSchemaCache(&mock, 10)
	for x := 0; x < 5; x++ {
		schema, err := schemaCache.GetSchema(0) // counter should only increment the first time
		if err != nil {
			t.Errorf("got unexpected error from GetSchema: %v", schema)
		}
		if !reflect.DeepEqual(schema, happySchema) {
			t.Errorf("expected schema (%v), got schema (%v)", happySchema, schema)
		}
	}
	if mock.counter != 1 {
		t.Errorf("expected cache hits but retrieved schema %d times", mock.counter)
	}
}
