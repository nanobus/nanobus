package widl_test

import (
	"os"
	"testing"

	"github.com/nanobus/nanobus/spec/widl"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	schemaBytes, err := os.ReadFile("../../example/customers/spec.widl")
	if err != nil {
		t.FailNow()
	}

	_, err = widl.Parse(schemaBytes)
	require.NoError(t, err)
	// jsonBytes, _ := json.MarshalIndent(ns, "", "  ")
	// fmt.Println(string(jsonBytes))
	// t.FailNow()
}
