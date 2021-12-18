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

	loader := func(with interface{}, resolver resolve.ResolveAs) (*compute.Compute, error) {
		return nil, nil
	}
	namedLoader := func() (string, compute.Loader) {
		return "test", loader
	}

	r.Register(namedLoader)

	assert.Equal(t, fmt.Sprintf("%v", compute.Loader(loader)), fmt.Sprintf("%p", r["test"]))
}
