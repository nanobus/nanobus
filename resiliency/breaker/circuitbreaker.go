package breaker

import (
	"time"

	"github.com/sony/gobreaker"

	"github.com/nanobus/nanobus/expr"
)

type CircuitBreaker struct {
	Name        string          `mapstructure:"name"`
	MaxRequests uint32          `mapstructure:"maxRequests"`
	Interval    time.Duration   `mapstructure:"interval"`
	Timeout     time.Duration   `mapstructure:"timeout"`
	Trip        *expr.ValueExpr `mapstructure:"trip"`
	breaker     *gobreaker.CircuitBreaker
}

func (c *CircuitBreaker) Initialize() {
	var tripFn func(counts gobreaker.Counts) bool = nil

	if c.Trip != nil {
		tripFn = func(counts gobreaker.Counts) bool {
			result, err := c.Trip.Eval(map[string]interface{}{
				"requests":             counts.Requests,
				"totalSuccesses":       counts.TotalSuccesses,
				"totalFailures":        counts.TotalFailures,
				"consecutiveSuccesses": counts.ConsecutiveSuccesses,
				"consecutiveFailures":  counts.ConsecutiveFailures,
			})
			if err != nil {
				// We cannot assume it is safe to trip if the eval
				// returns an error
				return false
			}
			if boolResult, ok := result.(bool); ok {
				return boolResult
			}

			return false
		}
	}

	c.breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        c.Name,
		MaxRequests: c.MaxRequests,
		Interval:    c.Interval,
		Timeout:     c.Timeout,
		ReadyToTrip: tripFn,
	})
}

func (c *CircuitBreaker) Execute(oper func() error) error {
	_, err := c.breaker.Execute(func() (interface{}, error) {
		err := oper()

		return nil, err
	})

	return err
}
