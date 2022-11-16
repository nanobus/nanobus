package dapr

import (
	"github.com/nanobus/nanobus/pkg/actions"
)

var All = []actions.NamedLoader{
	Publish,
	DeleteState,
	GetState,
	SetState,
	InvokeBinding,
}
