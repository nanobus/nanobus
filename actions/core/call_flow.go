package core

import (
	"context"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
)

type CallFlowConfig struct {
	// Name is the name of the flow to call.
	Name string `mapstructure:"name"`
}

// Route is the NamedLoader for the filter action.
func CallFlow() (string, actions.Loader) {
	return "call_flow", CallFlowLoader
}

func CallFlowLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	var c CallFlowConfig
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var processor Processor
	if err := resolve.Resolve(resolver,
		"system:processor", &processor); err != nil {
		return nil, err
	}

	return CallFlowAction(&c, processor), nil
}

func CallFlowAction(
	config *CallFlowConfig, processor Processor) actions.Action {
	return func(ctx context.Context, data actions.Data) (output interface{}, err error) {
		return processor.Flow(ctx, config.Name, data)
	}
}
