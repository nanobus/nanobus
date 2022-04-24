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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nanobus/nanobus/codec/confluentavro"
)

type cache struct {
	schema *confluentavro.Schema
}

func (c *cache) GetSchema(id int) (*confluentavro.Schema, error) {
	return c.schema, nil
}

var schema *confluentavro.Schema

func init() {
	schemaJSON := `{
		"type": "record",
		"name": "ExampleRecord",
		"namespace": "com.acme.messages",
		"fields": [
			{
				"name": "someProperty",
				"type": [
					"null",
					"string"
				]
			},
			{
				"name": "otherProperty",
				"type": {
					"type": "record",
					"name": "NestedRecord",
					"fields": [
						{
							"name": "nestedProperty",
							"type": "string"
						}
					]
				}
			}
		]
	}`
	var err error
	schema, err = confluentavro.ParseSchema(1, schemaJSON)
	if err != nil {
		panic(err)
	}
}

func TestEncodeDecode(t *testing.T) {
	record := map[string]interface{}{
		"someProperty": "foo",
		"otherProperty": map[string]interface{}{
			"nestedProperty": "bar",
		},
	}

	codec := confluentavro.NewCodec(&cache{schema: schema})
	encodedBytes, err := codec.Encode(record, schema.ID())
	require.Nil(t, err)
	require.Equal(t, []byte{0, 0, 0, 0, 1, 2, 6, 102, 111, 111, 6, 98, 97, 114}, encodedBytes)

	readI, _, err := codec.Decode(encodedBytes)
	require.NoError(t, err)
	read, ok := readI.(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, "foo", read["someProperty"])
	otherProperty, ok := read["otherProperty"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "bar", otherProperty["nestedProperty"])
}
