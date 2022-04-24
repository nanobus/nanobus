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

package function

import (
	"context"
)

type Function struct {
	Namespace string `json:"namespace" msgpack:"namespace"`
	Operation string `json:"operation" msgpack:"operation"`
}

type functionKey struct{}

func FromContext(ctx context.Context) Function {
	v := ctx.Value(functionKey{})
	if v == nil {
		return Function{}
	}
	c, _ := v.(Function)

	return c
}

func ToContext(ctx context.Context, function Function) context.Context {
	return context.WithValue(ctx, functionKey{}, function)
}
