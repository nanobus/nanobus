/*
 * Copyright 2022 The NanoBus Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package session

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"

	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
	"github.com/nanobus/nanobus/pkg/runtime"
	"github.com/nanobus/nanobus/pkg/transport/filter"
)

type Processor interface {
	LoadPipeline(pl *runtime.Pipeline) (runtime.Runnable, error)
	Pipeline(ctx context.Context, name string, data actions.Data) (interface{}, error)
	Provider(ctx context.Context, namespace, service, function string, data actions.Data) (interface{}, error)
	Event(ctx context.Context, name string, data actions.Data) (interface{}, error)
}

func SessionV1Loader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (filter.Filter, error) {
	var c SessionV1Config
	err := config.Decode(with, &c)
	if err != nil {
		return nil, err
	}

	var logger logr.Logger
	var processor runtime.Namespaces
	var developerMode bool
	if err := resolve.Resolve(resolver,
		"system:logger", &logger,
		"system:interfaces", &processor,
		"developerMode", &developerMode); err != nil {
		return nil, err
	}

	return Filter(logger, processor, &c, developerMode), nil
}

func Filter(log logr.Logger, processor runtime.Namespaces, config *SessionV1Config, developerMode bool) filter.Filter {
	return func(ctx context.Context, header filter.Header) (context.Context, error) {
		cookieHeader := header.Get("Cookie")
		hdr := http.Header{}
		hdr.Add("Cookie", cookieHeader)
		req := http.Request{Header: hdr}
		cookies := req.Cookies()

		var sid string
		for _, c := range cookies {
			if c.Name == "sid" {
				sid = c.Value
				break
			}
		}

		if sid == "" {
			return ctx, nil
		}

		result, _, err := processor.Invoke(ctx, config.Handler, actions.Data{
			"sid": sid,
		})
		if err != nil {
			return ctx, err
		}

		m, ok := result.(map[string]any)
		if !ok {
			return ctx, nil
		}

		var accessToken, tokenType string
		if accessTokenIface, ok := m["accessToken"]; ok {
			accessToken, _ = accessTokenIface.(string)
		}
		if tokenTypeIface, ok := m["tokenType"]; ok {
			tokenType, _ = tokenTypeIface.(string)
		}

		if developerMode {
			log.Info("Session debug info [TURN OFF FOR PRODUCTION]",
				"sid", sid,
				"token_type", tokenType,
				"access_token", accessToken)
		}

		header.Set("Authorization", tokenType+" "+accessToken)

		return ctx, nil
	}
}
