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

package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/resource"
	"github.com/nanobus/nanobus/spec"
)

type LoadConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `mapstructure:"resource" validate:"required"`
	// Namespace is the type namespace to load.
	Namespace string `mapstructure:"namespace" validate:"required"`
	// Type is the type name to load.
	Type string `mapstructure:"type" validate:"required"`
	// ID is the entity identifier expression.
	Key *expr.ValueExpr `mapstructure:"key" validate:"required"`
	// Preload lists the relationship to expand/load.
	Preload []Preload `mapstructure:"preload"`
	// NotFoundError is the error to return if the key is not found.
	NotFoundError string `mapstructure:"notFoundError"`
}

// Load is the NamedLoader for the invoke action.
func Load() (string, actions.Loader) {
	return "@postgres/load", LoadLoader
}

func LoadLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := LoadConfig{
		NotFoundError: "not_found",
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var namespaces spec.Namespaces
	var resources resource.Resources
	if err := resolve.Resolve(resolver,
		"spec:namespaces", &namespaces,
		"resource:lookup", &resources); err != nil {
		return nil, err
	}

	poolI, ok := resources[c.Resource]
	if !ok {
		return nil, fmt.Errorf("resource %q is not registered", c.Resource)
	}
	pool, ok := poolI.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("resource %q is not a *pgxpool.Pool", c.Resource)
	}

	ns, ok := namespaces[c.Namespace]
	if !ok {
		return nil, fmt.Errorf("namespace %q is not found", c.Namespace)
	}
	t, ok := ns.Type(c.Type)
	if !ok {
		return nil, fmt.Errorf("type %q is not found", c.Type)
	}

	return LoadAction(&c, t, ns, pool), nil
}

func LoadAction(
	config *LoadConfig,
	t *spec.Type,
	ns *spec.Namespace,
	pool *pgxpool.Pool) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		keyValue, err := config.Key.Eval(data)
		if err != nil {
			return nil, err
		}

		var result map[string]interface{}
		err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) (err error) {
			result, err = findById(ctx, conn, t, keyValue, config.Preload)
			return err
		})

		if result == nil && config.NotFoundError != "" {
			return nil, errorz.Return(config.NotFoundError, errorz.Metadata{
				"resource": config.Resource,
				"key":      keyValue,
			})
		}

		return result, err
	}
}
