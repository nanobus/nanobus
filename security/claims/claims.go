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

package claims

import (
	"context"
)

type Claims map[string]interface{}

func Combine(claimsList ...Claims) Claims {
	merged := make(Claims)

	for _, claims := range claimsList {
		if claims == nil {
			continue
		}

		for k, v := range claims {
			merged[k] = v
		}
	}

	return merged
}

type claimsKey struct{}

func FromContext(ctx context.Context) Claims {
	v := ctx.Value(claimsKey{})
	if v == nil {
		return Claims{}
	}
	c, _ := v.(Claims)
	if c == nil {
		return Claims{}
	}

	return c
}

func ToContext(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}
