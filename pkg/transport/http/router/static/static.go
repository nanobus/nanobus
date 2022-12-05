/*
 * Copyright 2022 The NanoBus Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

//go:generate apex generate
package static

import (
	"context"
	"net/http"
	"os"
	"sort"

	"github.com/go-logr/logr"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
	"github.com/nanobus/nanobus/pkg/transport/http/router"
)

func StaticV1Loader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (router.Router, error) {
	c := StaticV1Config{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var logger logr.Logger
	if err := resolve.Resolve(resolver,
		"system:logger", &logger); err != nil {
		return nil, err
	}

	return NewV1(logger, c), nil
}

func NewV1(log logr.Logger, config StaticV1Config) router.Router {
	return func(r *mux.Router, address string) error {
		sort.Slice(config.Paths, func(i, j int) bool {
			return len(config.Paths[i].Path) > len(config.Paths[j].Path)
		})
		for _, path := range config.Paths {
			log.Info("Serving static files",
				"dir", path.Dir,
				"path", path.Path,
				"strip", path.Strip)
			fs := http.FileServer(http.Dir(path.Dir))
			if path.Strip != nil {
				fs = http.StripPrefix(*path.Strip, fs)
			}
			r.PathPrefix(path.Path).Handler(handlers.LoggingHandler(os.Stdout, fs))
		}

		return nil
	}
}
