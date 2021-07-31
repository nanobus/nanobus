package runtime

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resiliency/breaker"
	"github.com/nanobus/nanobus/resiliency/retry"
	"github.com/nanobus/nanobus/resolve"
)

type Environment map[string]string

type Processor struct {
	config          *Configuration
	registry        actions.Registry
	resolver        resolve.DependencyResolver
	resolveAs       resolve.ResolveAs
	retries         map[string]*retry.Config
	circuitBreakers map[string]*breaker.CircuitBreaker
	services        Namespaces
	outbound        Namespaces
	inbound         Functions
}

type Namespaces map[string]Functions
type Functions map[string]Runnable

type Runnable struct {
	config *Pipeline
	steps  []step
}

type step struct {
	config         *Step
	action         actions.Action
	timeout        *time.Duration
	retry          *retry.Config
	circuitBreaker *breaker.CircuitBreaker
}

func New(configuration *Configuration, registry actions.Registry, resolver resolve.DependencyResolver) (*Processor, error) {
	retries := make(map[string]*retry.Config, len(configuration.Resiliency.Retries))
	for name, retryMap := range configuration.Resiliency.Retries {
		retryConfig, err := retry.DecodeConfig(retryMap)
		if err != nil {
			return nil, err
		}
		retries[name] = &retryConfig
	}

	circuitBreakers := make(map[string]*breaker.CircuitBreaker, len(configuration.Resiliency.CircuitBreakers))
	for name, circuitBreaker := range configuration.Resiliency.CircuitBreakers {
		var cb breaker.CircuitBreaker
		if err := config.Decode(circuitBreaker, &cb); err != nil {
			return nil, err
		}
		cb.Initialize(name)
		circuitBreakers[name] = &cb
	}

	p := Processor{
		config:          configuration,
		retries:         retries,
		circuitBreakers: circuitBreakers,
		registry:        registry,
	}

	p.resolver = func(name string) (interface{}, bool) {
		if name == "system:processor" {
			return p, true
		}
		return resolver(name)
	}

	p.resolveAs = resolve.ToResolveAs(p.resolver)

	if err := p.initialize(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Processor) Service(ctx context.Context, namespace, service, function string, data actions.Data) (interface{}, bool, error) {
	s, ok := p.services[namespace+"."+service]
	if !ok {
		return nil, false, nil
	}

	pl, ok := s[function]
	if !ok {
		return nil, false, nil
	}

	output, err := pl.Run(ctx, data)
	if err == nil && output == nil {
		output = data["input"]
	}
	return output, true, err
}

func (p *Processor) Outbound(ctx context.Context, namespace, service, function string, data actions.Data) (interface{}, error) {
	nss := namespace + "." + service
	s, ok := p.outbound[nss]
	if !ok {
		return nil, fmt.Errorf("unknown outbound service %q", nss)
	}

	pl, ok := s[function]
	if !ok {
		return nil, fmt.Errorf("unknown outbound function %q in service %q", function, nss)
	}

	return pl.Run(ctx, data)
}

func (p *Processor) Inbound(ctx context.Context, function string, data actions.Data) (interface{}, error) {
	pl, ok := p.inbound[function]
	if !ok {
		return nil, fmt.Errorf("unknown inbound function %q", function)
	}

	return pl.Run(ctx, data)
}

func (p *Processor) initialize() (err error) {
	if p.services, err = p.loadServices(p.config.Services); err != nil {
		return err
	}
	if p.outbound, err = p.loadServices(p.config.Outbound); err != nil {
		return err
	}
	if p.inbound, err = p.loadFunctionPipelines(p.config.Inbound); err != nil {
		return err
	}

	return nil
}

func (p *Processor) loadServices(services Services) (s Namespaces, err error) {
	s = make(Namespaces, len(services))
	for ns, fns := range services {
		if s[ns], err = p.loadFunctionPipelines(fns); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (p *Processor) loadFunctionPipelines(fpl FunctionPipelines) (Functions, error) {
	runnables := make(Functions, len(fpl))
	for name, pipeline := range fpl {
		pl, err := p.LoadPipeline(&pipeline)
		if err != nil {
			return nil, err
		}
		runnables[name] = *pl
	}

	return runnables, nil
}

func (p *Processor) LoadPipeline(pl *Pipeline) (*Runnable, error) {
	steps := make([]step, len(pl.Actions))
	for i, s := range pl.Actions {
		step, err := p.loadStep(&s)
		if err != nil {
			return nil, err
		}
		steps[i] = *step
	}

	return &Runnable{
		config: pl,
		steps:  steps,
	}, nil
}

func (p *Processor) loadStep(s *Step) (*step, error) {
	loader, ok := p.registry[s.Name]
	if !ok {
		return nil, fmt.Errorf("unregistered action %q", s.Name)
	}

	action, err := loader(s.With, p.resolveAs)
	if err != nil {
		return nil, err
	}

	var retry *retry.Config
	if s.Retry != "" {
		var ok bool
		retry, ok = p.retries[s.Retry]
		if !ok {
			return nil, fmt.Errorf("retry policy %q is not defined", s.Retry)
		}
	}

	var circuitBreaker *breaker.CircuitBreaker
	if s.CircuitBreaker != "" {
		var ok bool
		circuitBreaker, ok = p.circuitBreakers[s.CircuitBreaker]
		if !ok {
			return nil, fmt.Errorf("circuit breaker policy %q is not defined", s.CircuitBreaker)
		}
	}

	var timeout *time.Duration
	if s.Timeout != "" {
		to, err := time.ParseDuration(s.Timeout)
		if err != nil {
			return nil, err
		}
		timeout = &to
	}

	return &step{
		config:         s,
		action:         action,
		timeout:        timeout,
		retry:          retry,
		circuitBreaker: circuitBreaker,
	}, nil
}

func (r *Runnable) Run(ctx context.Context, data actions.Data) (interface{}, error) {
	var output interface{}
	var err error
	for _, s := range r.steps {
		rp := ResiliencyPolicy{
			Name:           s.config.Summary,
			Timeout:        s.timeout,
			Retry:          s.retry,
			CircuitBreaker: s.circuitBreaker,
		}
		err = rp.Run(ctx, func(ctx context.Context) error {
			output, err = s.action(ctx, data)
			if errors.Is(err, actions.ErrStop) {
				return backoff.Permanent(err)
			}
			return err
		})
	}

	return output, nil
}

type ResiliencyPolicy struct {
	Name           string
	Timeout        *time.Duration          `mapstructure:"timeout"`
	Retry          *retry.Config           `mapstructure:"retry"`
	CircuitBreaker *breaker.CircuitBreaker `mapstructure:"circuitBreaker"`
}

type Operation func(ctx context.Context) error

func (p *ResiliencyPolicy) Run(ctx context.Context, oper Operation) error {
	operation := oper
	if p.Timeout != nil {
		// Handle timeout
		operation = func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, *p.Timeout)
			defer cancel()

			return oper(ctx)
		}
	}

	var call func() error
	if p.Retry == nil {
		call = func() error {
			return operation(ctx)
		}
	} else {
		// Use retry/back off
		b := p.Retry.NewBackOffWithContext(ctx)
		call = func() error {
			return retry.NotifyRecover(func() error {
				return operation(ctx)
			}, b, func(_ error, _ time.Duration) {
				log.Printf("Error processing operation %s. Retrying...", p.Name)
			}, func() {
				log.Printf("Recovered processing operation %s.", p.Name)
			})
		}
	}

	if p.CircuitBreaker != nil {
		return p.CircuitBreaker.Execute(call)
	}

	return call()
}
