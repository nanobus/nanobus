package gorm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/resource"
	"github.com/nanobus/nanobus/spec"
	"github.com/nanobus/nanobus/stream"
	"gorm.io/gorm"
)

type FindConfig struct {
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
func Find() (string, actions.Loader) {
	return "@gorm/find", FindLoader
}

func FindLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := FindConfig{
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
	pool, ok := poolI.(*gorm.DB)
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

func annotationValue(a spec.Annotator, annotation, argument, defaultValue string) string {
	if av, ok := a.Annotation(annotation); ok {
		if arg, ok := av.Argument(argument); ok {
			return fmt.Sprintf("%v", arg.Value)
		}
	}
	return defaultValue
}

func FindAction(
	config *FindConfig,
	t *spec.Type,
	ns *spec.Namespace,
	db *gorm.DB) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		s, ok := stream.FromContext(ctx)
		if !ok {
			return nil, errors.New("stream not in context")
		}

		//table := annotationValue(t, "entity", "table", t.Name)

		p := NewProcessor(db.NamingStrategy)
		if err := p.ConvertTypes(ns.Types); err != nil {
			return nil, err
		}

		// schemas := make(map[string]*schema.Schema)
		// for _, d := range ns.Types {
		// 	TypeToSchema(schemas, d, db.NamingStrategy)
		// }

		pair, err := p.TypeToSchema(t)
		if err != nil {
			return nil, err
		}

		// db.NamingStrategy
		// db.Statement.Schema

		tx := db.Table(pair.S.Table)
		tx.Statement.Schema = pair.S
		tx = tx.Preload("address")
		tx.Statement.Schema = pair.S

		var results []map[string]interface{}
		tx = tx.Find(&results)
		if tx.Error != nil {
			return nil, err
		}

		for _, result := range results {
			fmt.Println(result)
			if err = s.SendData(result); err != nil {
				return nil, err
			}
		}

		// rows, err := db.Preload("Address").Table(schema.Table).Find(&results)
		// if err != nil {
		// 	return nil, err
		// }
		// defer rows.Close()

		// columns, err := rows.Columns()
		// if err != nil {
		// 	return nil, err
		// }

		// fields := make([]*spec.Field, len(columns))
		// types := make([]reflect.Type, len(columns))
		// for i, col := range columns {
		// 	for _, field := range t.Fields {
		// 		colname := annotationValue(field, "column", "name", field.Name)
		// 		if colname == col {
		// 			fields[i] = field
		// 			types[i] = reflectType(field.Type)
		// 			break
		// 		}
		// 	}
		// }

		// for rows.Next() {
		// 	item := make(map[string]interface{}, len(columns))
		// 	values := make([]interface{}, len(columns))
		// 	for idx, t := range types {
		// 		if t != nil {
		// 			values[idx] = reflect.New(reflect.PtrTo(t)).Interface()
		// 		}
		// 	}
		// 	rows.Scan(values...)

		// 	for i, field := range fields {
		// 		if field == nil {
		// 			continue
		// 		}
		// 		item[field.Name] = values[i]
		// 	}

		// 	if err = s.SendData(item); err != nil {
		// 		return nil, err
		// 	}
		// }

		return nil, err
	}
}

func reflectType(t *spec.TypeRef) reflect.Type {
	switch t.Kind {
	case spec.KindString:
		return reflect.TypeOf("")
	case spec.KindU64:
		return reflect.TypeOf(uint64(1234))
	case spec.KindOptional:
		return reflectType(t.OptionalType)
	case spec.KindType:
		return mapType
	}
	fmt.Println(t.Kind, t)
	return nil
}

func getOne(ctx context.Context, conn *pgxpool.Conn, config *FindConfig, t *spec.Type, idColumn string, idValue interface{}, toPreload []Preload) (interface{}, error) {
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
