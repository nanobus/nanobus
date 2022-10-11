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

package gorm

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/resource"
	"github.com/nanobus/nanobus/spec"
	"github.com/nanobus/nanobus/stream"
)

type FindAllConfig struct {
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
}

// Load is the NamedLoader for the invoke action.
func FindAll() (string, actions.Loader) {
	return "@gorm/find_all", FindAllLoader
}

func FindAllLoader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := FindAllConfig{}
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

	return FindAllAction(&c, t, ns, pool), nil
}

func FindAllAction(
	config *FindAllConfig,
	t *spec.Type,
	ns *spec.Namespace,
	db *gorm.DB) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		s, ok := stream.SinkFromContext(ctx)
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
			if err = s.Next(result, nil); err != nil {
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
