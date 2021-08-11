package core

import (
	"net/http"

	"github.com/nanobus/nanobus/actions"
)

var All = []actions.NamedLoader{
	Assign,
	Decode,
	Filter,
	Invoke,
	Log,
	Route,
}

// Common dependencies

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
