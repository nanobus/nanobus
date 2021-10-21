package claims

import (
	"context"
)

type Claims map[string]interface{}

func (c Claims) Combine(other Claims) Claims {
	if other == nil {
		return c
	}
	if c == nil {
		return other
	}

	merged := make(Claims, len(c)+len(other))
	for k, v := range c {
		merged[k] = v
	}

	for k, v := range other {
		merged[k] = v
	}

	return merged
}

type claimsKey struct{}

func FromContext(ctx context.Context) Claims {
	v := ctx.Value(claimsKey{})
	if v == nil {
		return nil
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
