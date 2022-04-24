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

				done := make(chan error, 1)
				go func() {
					done <- operCopy(ctx)
				}()

				select {
				case err := <-done:
					return err
				case <-ctx.Done():
					return ctx.Err()
				}
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
		err := retry.NotifyRecover(func() error {
			return operation(ctx)
		}, b, func(err error, _ time.Duration) {
			log.Error(err, "Error processing operation. Retrying...", "operation", operationName)
		}, func() {
			log.Info("Recovered processing operation.", "operation", operationName)
		})
		if err != nil {
			if perr, ok := err.(*backoff.PermanentError); ok {
				err = perr.Err
			}
		}
		return err
	}
}
