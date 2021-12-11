package breaker_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resiliency/breaker"
)

func TestCircuitBreaker(t *testing.T) {
	var trip expr.ValueExpr
	err := trip.DecodeString("consecutiveFailures > 2")
	require.NoError(t, err)
	cb := breaker.CircuitBreaker{
		Name: "test",
		Trip: &trip,
	}
	cb.Initialize()
	for i := 0; i < 3; i++ {
		cb.Execute(func() error {
			return errors.New("test")
		})
	}
	err = cb.Execute(func() error {
		return nil
	})
	assert.EqualError(t, err, "circuit breaker is open")
}
