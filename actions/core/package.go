package core

import (
	"context"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/runtime"
)

var All = []actions.NamedLoader{
	Assign,
	Authorize,
	CallFlow,
	Decode,
	Filter,
	HTTP,
	Invoke,
	Log,
	Route,
}

type Processor interface {
	LoadPipeline(pl *runtime.Pipeline) (runtime.Runnable, error)
	Flow(ctx context.Context, name string, data actions.Data) (interface{}, error)
}
