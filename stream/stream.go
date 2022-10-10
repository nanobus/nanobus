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

type Source interface {
	Next(data any, md *metadata.MD) error
}

type sourceKey struct{}

// SourceNewContext creates a new context with incoming `s` attached.
func SourceNewContext(ctx context.Context, s Source) context.Context {
	return context.WithValue(ctx, sourceKey{}, s)
}

func SourceFromContext(ctx context.Context) (s Source, ok bool) {
	s, ok = ctx.Value(sourceKey{}).(Source)
	return
}

type Sink interface {
	Next(data any, md metadata.MD) error
	Complete()
	Error(err error)
}

type sinkKey struct{}

// SinkNewContext creates a new context with incoming `s` attached.
func SinkNewContext(ctx context.Context, s Sink) context.Context {
	return context.WithValue(ctx, sinkKey{}, s)
}

func SinkFromContext(ctx context.Context) (s Sink, ok bool) {
	s, ok = ctx.Value(sinkKey{}).(Sink)
	return
}
