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

package paseto

import (
	"context"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-logr/logr"

	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/security/claims"
	"github.com/nanobus/nanobus/transport/filter"
)

type Config struct {
	Audience     string `mapstructure:"audience"`
	Issuer       string `mapstructure:"issuer"`
	PublicKeyHex string `mapstructure:"publicKeyHex" validate:"required"`
}

type Settings struct {
}

// Paseto is the NamedLoader for the Paseto filter.
func Paseto() (string, filter.Loader) {
	return "paseto", Loader
}

func Loader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (filter.Filter, error) {
	var c Config
	err := config.Decode(with, &c)
	if err != nil {
		return nil, err
	}

	var logger logr.Logger
	if err := resolve.Resolve(resolver,
		"system:logger", &logger); err != nil {
		return nil, err
	}

	parser := paseto.NewParser()
	if c.Audience != "" {
		parser.AddRule(paseto.ForAudience(c.Audience))
	}
	if c.Issuer != "" {
		parser.AddRule(paseto.IssuedBy(c.Issuer))
	}
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(c.PublicKeyHex)
	if err != nil {
		// panic or deal with error of invalid key
		return nil, err
	}

	return Filter(logger, parser, publicKey), nil
}

func Filter(log logr.Logger, parser paseto.Parser, publicKey paseto.V4AsymmetricPublicKey) filter.Filter {
	return func(ctx context.Context, header filter.Header) (context.Context, error) {
		authorization := header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			return ctx, nil
		}

		tokenString := authorization[7:]

		parsedToken, err := parser.ParseV4Public(publicKey, tokenString, nil)
		if err != nil {
			// deal with error of token which failed to be validated, or cryptographically verified
			return nil, err
		}

		tokenClaims := parsedToken.Claims()
		if tokenClaims != nil {
			ctx = claims.ToContext(ctx, claims.Claims(tokenClaims))
		}

		return ctx, nil
	}
}
