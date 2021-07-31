package core

import (
	"net/http"

	"github.com/nanobus/nanobus/actions"
)

var All = []actions.NamedLoader{
	Filter,
	Invoke,
	Log,
	Route,
	Assign,
}

// Common dependencies

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
