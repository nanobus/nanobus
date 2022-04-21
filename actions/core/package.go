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
	CallProvider,
	Decode,
	Filter,
	HTTP,
	Invoke,
	JMESPath,
	JQ,
	Log,
	Route,
}

type Processor interface {
	LoadPipeline(pl *runtime.Pipeline) (runtime.Runnable, error)
	Flow(ctx context.Context, name string, data actions.Data) (interface{}, error)
	Provider(ctx context.Context, namespace, service, function string, data actions.Data) (interface{}, error)
	Event(ctx context.Context, name string, data actions.Data) (interface{}, error)
}
