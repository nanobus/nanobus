package postgres

import (
	"context"
	"fmt"
	"strings"

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
	Resource string `mapstructure:"resource"`
	// Namespace is the type namespace to load.
	Namespace string `mapstructure:"namespace"`
	// Type is the type name to load.
	Type string `mapstructure:"type"`
	// ID is the entity identifier expression.
	ID *expr.ValueExpr `mapstructure:"id"`
	// Preload lists the relationship to expand/load.
	Preload []Preload `mapstructure:"preload"`
	// NotFoundError is the error to return if the key is not found.
	NotFoundError string `mapstructure:"notFoundError"`
}

type Preload struct {
	Field   string    `mapstructure:"field"`
	Preload []Preload `mapstructure:"preload"`
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

func annotationValue(a spec.Annotator, annotation, argument, defaultValue string) string {
	if av, ok := a.Annotation(annotation); ok {
		if arg, ok := av.Argument(argument); ok {
			return fmt.Sprintf("%v", arg.Value)
		}
	}
	return defaultValue
}

func LoadAction(
	config *LoadConfig,
	t *spec.Type,
	ns *spec.Namespace,
	pool *pgxpool.Pool) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		idColumn := keyColumn(t)
		idValue, err := config.ID.Eval(data)
		if err != nil {
			return nil, err
		}

		var result interface{}
		err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) (err error) {
			result, err = getOne(ctx, conn, config, t, idColumn, idValue, config.Preload)
			return err
		})

		if result == nil {
			return nil, errorz.Return(config.NotFoundError, errorz.Metadata{
				"resource": config,
				"key":      idValue,
			})
		}

		return result, err
	}
}

func getOne(ctx context.Context, conn *pgxpool.Conn, config *LoadConfig, t *spec.Type, idColumn string, idValue interface{}, toPreload []Preload) (interface{}, error) {
	keyCol := keyColumn(t)
	sql := generateTableSQL(t)
	rows, err := conn.Query(ctx, sql+" WHERE "+idColumn+"=$1", idValue)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		record := make(map[string]interface{})
		values, err := rows.Values()
		rows.Close()
		if err != nil {
			return nil, err
		}
		for i, v := range values {
			record[t.Fields[i].Name] = v
		}

		for _, preload := range toPreload {
			ex, ok := t.Field(preload.Field)
			if !ok {
				return nil, fmt.Errorf("%s is not a field of %s", preload.Field, t.Name)
			}
			fk := annotationValue(ex, "hasOne", "key", "")
			if fk == "" {
				return nil, fmt.Errorf("hasOne is not specified on %s", ex.Name)
			}

			res, err := getOne(ctx, conn, config, ex.Type.Type, fk, record[keyCol], preload.Preload)
			if err != nil {
				return nil, err
			}

			record[preload.Field] = res
		}

		return record, nil
	}

	rows.Close()

	return nil, nil
}

func keyColumn(t *spec.Type) string {
	for _, f := range t.Fields {
		if _, ok := f.Annotation("key"); ok {
			return annotationValue(t, "column", "name", f.Name)
		}
	}
	return ""
}

func generateTableSQL(t *spec.Type) string {
	var buf strings.Builder

	buf.WriteString("SELECT ")
	for i, f := range t.Fields {
		column := annotationValue(f, "column", "name", "")
		if column == "" {
			continue
		}
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(column)
	}
	buf.WriteString(" FROM ")
	table := annotationValue(t, "entity", "table", t.Name)
	buf.WriteString(table)

	return buf.String()
}
