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

package stream

import (
	"context"

	"github.com/nanobus/nanobus/channel/metadata"
)

type Stream interface {
	Metadata() metadata.MD
	RecvData(dst interface{}) error
	SendHeaders(md metadata.MD, end ...bool) error
	SendData(v interface{}, end ...bool) error
	SendUnary(md metadata.MD, v interface{}) error
	SendRequest(path string, v interface{}) error
	SendReply(v interface{}) error
	SendError(err error) error
}

type streamKey struct{}

// NewContext creates a new context with incoming `s` attached.
func NewContext(ctx context.Context, s Stream) context.Context {
	return context.WithValue(ctx, streamKey{}, s)
}

func FromContext(ctx context.Context) (s Stream, ok bool) {
	s, ok = ctx.Value(streamKey{}).(Stream)
	return
}
