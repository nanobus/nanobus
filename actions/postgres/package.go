package postgres

import (
	"github.com/nanobus/nanobus/actions"
)

var All = []actions.NamedLoader{
	Load,
	Find,
	Query,
	Exec,
	Test,
}
