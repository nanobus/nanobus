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

package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nanobus/nanobus/codec/json"
)

func TestCodec(t *testing.T) {
	name, auto, loader := json.JSON()
	assert.Equal(t, "json", name)
	assert.True(t, auto)
	c, err := loader(nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "application/json", c.ContentType())
	data := map[string]interface{}{
		"int":    int64(1234),
		"string": "1234",
	}
	encoded, err := c.Encode(data)
	require.NoError(t, err)
	_, _, err = c.Decode([]byte(`bad data`))
	assert.Error(t, err)
	decoded, _, err := c.Decode(encoded)
	require.NoError(t, err)
	assert.Equal(t, data, decoded)
}
