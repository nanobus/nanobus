package actions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/resolve"
)

func TestRegistry(t *testing.T) {
	r := actions.Registry{}

	loader := func(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
		return nil, nil
	}
	namedLoader := func() (string, actions.Loader) {
		return "test", loader
	}

	r.Register(namedLoader)

	assert.Equal(t, fmt.Sprintf("%v", actions.Loader(loader)), fmt.Sprintf("%p", r["test"]))
}

func TestClone(t *testing.T) {
	data := actions.Data{
		"one": 1234,
	}
	clone := data.Clone()
	assert.Equal(t, data, clone)
	assert.NotSame(t, data, clone)
}

func TestStop(t *testing.T) {
	assert.Equal(t, actions.ErrStop, actions.Stop())
}
