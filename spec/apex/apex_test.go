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

package apex_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nanobus/nanobus/spec/apex"
)

func TestParse(t *testing.T) {
	expectedBytes, err := os.ReadFile("testdata/expected.json")
	if err != nil {
		t.FailNow()
	}

	name, loader := apex.Apex()
	assert.Equal(t, "apex", name)
	namespaces, err := loader(map[string]interface{}{
		"filename": "testdata/spec.apexlang",
	})
	require.NoError(t, err)
	require.Len(t, namespaces, 1)

	actualBytes, err := json.MarshalIndent(namespaces[0], "", "  ")
	require.NoError(t, err)
	fmt.Println(string(actualBytes))

	var expected, actual interface{}
	require.NoError(t, json.Unmarshal(expectedBytes, &expected))
	require.NoError(t, json.Unmarshal(actualBytes, &actual))

	assert.Equal(t, expected, actual)
}
