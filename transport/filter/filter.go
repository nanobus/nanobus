package filter

import (
	"context"
	"net/http"
)

type Filter func(ctx context.Context, req *http.Request) (context.Context, error)
