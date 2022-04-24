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
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanobus/nanobus/spec"
)

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

func annotationValue(a spec.Annotator, annotation, argument, defaultValue string) string {
	if av, ok := a.Annotation(annotation); ok {
		if arg, ok := av.Argument(argument); ok {
			return fmt.Sprintf("%v", arg.Value)
		}
	}
	return defaultValue
}
