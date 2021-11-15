package filter_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/transport/filter"
)

func TestRegistry(t *testing.T) {
	r := filter.Registry{}

	loader := func(with interface{}, resolver resolve.ResolveAs) (filter.Filter, error) {
		return nil, nil
	}
	namedLoader := func() (string, filter.Loader) {
		return "test", loader
	}

	r.Register(namedLoader)

	assert.Equal(t, fmt.Sprintf("%v", filter.Loader(loader)), fmt.Sprintf("%p", r["test"]))
}
