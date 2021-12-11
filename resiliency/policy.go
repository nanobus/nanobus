package resiliency

import (
	"context"
	"log"
	"time"

	"github.com/nanobus/nanobus/resiliency/breaker"
	"github.com/nanobus/nanobus/resiliency/retry"
)

type (
	Policy struct {
		operationName  string
		timeout        time.Duration           `mapstructure:"timeout"`
		retry          *retry.Config           `mapstructure:"retry"`
		circuitBreaker *breaker.CircuitBreaker `mapstructure:"circuitBreaker"`
	}

	Operation func(ctx context.Context) error
)

func NewPolicy(operationName string, t time.Duration, r *retry.Config, cb *breaker.CircuitBreaker) Policy {
	return Policy{
		operationName:  operationName,
		timeout:        t,
		retry:          r,
		circuitBreaker: cb,
	}
}

func (p *Policy) Run(ctx context.Context, oper Operation) error {
	operation := oper
	if p.timeout > 0 {
		// Handle timeout
		operation = func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, p.timeout)
			defer cancel()

			return oper(ctx)
		}
	}

	var call func() error
	if p.retry == nil {
		call = func() error {
			return operation(ctx)
		}
	} else {
		// Use retry/back off
		b := p.retry.NewBackOffWithContext(ctx)
		call = func() error {
			return retry.NotifyRecover(func() error {
				return operation(ctx)
			}, b, func(_ error, _ time.Duration) {
				log.Printf("Error processing operation %s. Retrying...", p.operationName)
			}, func() {
				log.Printf("Recovered processing operation %s.", p.operationName)
			})
		}
	}

	if p.circuitBreaker != nil {
		return p.circuitBreaker.Execute(call)
	}

	return call()
}
