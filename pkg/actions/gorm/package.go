package gorm

import (
	"github.com/nanobus/nanobus/pkg/actions"
)

var All = []actions.NamedLoader{
	Find,
	FindAll,
}
