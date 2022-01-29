package dapr

import (
	"github.com/nanobus/nanobus/actions"
)

var All = []actions.NamedLoader{
	Inoke,
	InvokeActor,
	InvokeBinding,
	SetState,
	GetState,
	DeleteState,
	PublishMessage,
	SQLExec,
}
