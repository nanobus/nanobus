/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httpresponse

import (
	"context"
	"net/http"
)

type Response struct {
	Status int
	Header http.Header
}

func New() *Response {
	return &Response{
		Status: http.StatusOK,
		Header: http.Header{},
	}
}

type responseKey struct{}

// NewContext creates a new context with incoming `resp` attached.
func NewContext(ctx context.Context, resp *Response) context.Context {
	return context.WithValue(ctx, responseKey{}, resp)
}

func FromContext(ctx context.Context) *Response {
	return ctx.Value(responseKey{}).(*Response)
}
