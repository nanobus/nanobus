package resiliency

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-logr/logr"

	"github.com/nanobus/nanobus/resiliency/breaker"
	"github.com/nanobus/nanobus/resiliency/retry"
)

type (
	// Operation represents a function to invoke with resiliency policies applied.
	Operation func(ctx context.Context) error

	// Runner represents a function to invoke `oper` with resiliency policies applied.
	Runner func(ctx context.Context, oper Operation) error
)

// Policy returns a policy runner that encapsulates the configured
// resiliency policies in a simple execution wrapper.
func Policy(log logr.Logger, operationName string, t time.Duration, r *retry.Config, cb *breaker.CircuitBreaker) Runner {
	return func(ctx context.Context, oper Operation) error {
		operation := oper
		if t > 0 {
			// Handle timeout
			operCopy := operation
			operation = func(ctx context.Context) error {
				ctx, cancel := context.WithTimeout(ctx, t)
				defer cancel()

				return operCopy(ctx)
			}
		}

		if cb != nil {
			operCopy := operation
			operation = func(ctx context.Context) error {
				err := cb.Execute(func() error {
					return operCopy(ctx)
				})
				if r != nil && breaker.IsErrorPermanent(err) {
					// Break out of retry.
					err = backoff.Permanent(err)
				}
				return err
			}
		}

		if r == nil {
			return operation(ctx)
		}

		// Use retry/back off
		b := r.NewBackOffWithContext(ctx)
		return retry.NotifyRecover(func() error {
			return operation(ctx)
		}, b, func(_ error, _ time.Duration) {
			log.Info("Error processing operation. Retrying...", "operation", operationName)
		}, func() {
			log.Info("Recovered processing operation.", "operation", operationName)
		})
	}
}
