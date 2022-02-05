package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/resource"
	"github.com/nanobus/nanobus/spec"
)

type FindConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `mapstructure:"resource"`
	// Namespace is the type namespace to load.
	Namespace string `mapstructure:"namespace"`
	// Type is the type name to load.
	Type string `mapstructure:"type"`
	// Preload lists the relationship to expand/load.
	Preload []Preload `mapstructure:"preload"`
	// Where list the parts of the where clause.
	Where []Where `mapstructure:"where"`
}

// Find is the NamedLoader for the invoke action.
func Find() (string, actions.Loader) {
	return "@postgres/find", FindLoader
}

func FindLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := FindConfig{}
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

	return FindAction(&c, t, ns, pool), nil
}

func FindAction(
	config *FindConfig,
	t *spec.Type,
	ns *spec.Namespace,
	pool *pgxpool.Pool) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		var result interface{}
		err := pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) (err error) {
			result, err = getMany(ctx, conn, t, data, config.Where, config.Preload)
			return err
		})

		return result, err
	}
}
