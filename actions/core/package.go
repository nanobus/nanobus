package core

import (
	"github.com/nanobus/nanobus/actions"
)

var All = []actions.NamedLoader{
	Assign,
	Authorize,
	Decode,
	Filter,
	HTTP,
	Invoke,
	Log,
	Route,
}
