package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/runtime"
)

// SelectionMode indicates how many routes can be selected.
type SelectionMode int

const (
	// Single indicates only one route can be selected.
	Single SelectionMode = iota
	// Multi indicates many routes can be selected.
	Multi
)

type RouteConfig struct {
	// Selection defines the selection mode: single or multi.
	Selection SelectionMode `mapstructure:"selection"`
	// Routes are the possible runnable routes which conditions for selection.
	Routes []RouteCondition `mapstructure:"routes"`
}

type RouteCondition struct {
	// Summary if the overall summary of this route.
	Summary string
	// When is the predicate expression for filtering.
	When *expr.ValueExpr `mapstructure:"when"`
	// Then is the steps to process.
	Then []runtime.Step `mapstructure:"then"`

	runnable *runtime.Runnable
}

// Route is the NamedLoader for the filter action.
func Route() (string, actions.Loader) {
	return "route", RouteLoader
}

func RouteLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c RouteConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var processor runtime.Processor
	if err := resolve.Resolve(resolver,
		"system:processor", &processor); err != nil {
		return nil, err
	}

	for i := range c.Routes {
		r := &c.Routes[i]
		runnable, err := processor.LoadPipeline(&runtime.Pipeline{
			Summary: r.Summary,
			Actions: r.Then,
		})
		if err != nil {
			return nil, err
		}
		r.runnable = runnable
	}

	return RouteAction(&c), nil
}

func RouteAction(
	config *RouteConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		for i := range config.Routes {
			r := &config.Routes[i]
			resultInt, err := r.When.Eval(data)
			if err != nil {
				return nil, err
			}

			result, ok := resultInt.(bool)
			if !ok {
				return nil, fmt.Errorf("expression %q did not evaluate a boolean", r.When.Expr())
			}

			if !result {
				continue
			}

			output, err := r.runnable.Run(ctx, data)
			if config.Selection == Single || err != nil {
				return output, err
			}
		}

		return nil, nil
	}
}

// DecodeString handles converting a string value to SelectionMode.
func (sm *SelectionMode) DecodeString(value string) error {
	switch strings.ToLower(value) {
	case "single":
		*sm = Single
	case "multi":
		*sm = Multi
	default:
		return fmt.Errorf("unexpected selection mode: %s", value)
	}

	return nil
}
