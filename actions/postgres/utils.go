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
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/spec"
)

type Preload struct {
	Field   string    `mapstructure:"field"`
	Preload []Preload `mapstructure:"preload"`
}

type Where struct {
	Query string          `mapstructure:"query"`
	Value *expr.ValueExpr `mapstructure:"value"`
}

func annotationValue(a spec.Annotator, annotation, argument, defaultValue string) string {
	if av, ok := a.Annotation(annotation); ok {
		if arg, ok := av.Argument(argument); ok {
			return fmt.Sprintf("%v", arg.Value)
		}
	}
	return defaultValue
}

func findById(ctx context.Context, conn *pgxpool.Conn, t *spec.Type, idValue interface{}, toPreload []Preload) (map[string]interface{}, error) {
	idColumn := keyColumn(t)
	sql := generateTableSQL(t) + " WHERE " + idColumn + " = $1"
	rows, err := conn.Query(ctx, sql, idValue)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		record := make(map[string]interface{})
		values, err := rows.Values()
		if err != nil {
			rows.Close()
			return nil, err
		}
		for i, v := range values {
			v = normalizeValue(v)
			record[t.Fields[i].Name] = v
		}

		rows.Close()

		for _, preload := range toPreload {
			ex, ok := t.Field(preload.Field)
			if !ok {
				return nil, fmt.Errorf("%s is not a field of %s", preload.Field, t.Name)
			}

			var res interface{}
			if _, ok := ex.Annotation("hasOne"); ok {
				res, err = findById(ctx, conn, ex.Type.Type, record[preload.Field], preload.Preload)
				if err != nil {
					return nil, err
				}
			} else if hasMany, ok := ex.Annotation("hasMany"); ok {
				if key, ok := hasMany.Argument("key"); ok {
					keyName := keyField(t)
					res, err = join(ctx, conn, ex.Type.ItemType.Type,
						key.ValueString()+" = $1", []interface{}{record[keyName]},
						preload.Preload)
					if err != nil {
						return nil, err
					}
				}
			}

			record[preload.Field] = res
		}

		return record, nil
	}

	return nil, nil
}

func findOne(ctx context.Context, conn *pgxpool.Conn, t *spec.Type, input map[string]interface{}, where []Where, toPreload []Preload) (map[string]interface{}, error) {
	sql := generateTableSQL(t)
	args := []interface{}{}
	if len(where) > 0 {
		dollarIndex := 1
		for i, part := range where {
			val, err := part.Value.Eval(input)
			if err != nil {
				return nil, err
			}
			if isNil(val) {
				continue
			}
			if i > 0 {
				sql += " AND "
			} else {
				sql += " WHERE "
			}
			query := part.Query
			for strings.Contains(query, "?") {
				query = strings.Replace(query, "?", fmt.Sprintf("$%d", dollarIndex), 1)
				dollarIndex++
			}
			sql += query
			args = append(args, val)
		}
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		record := make(map[string]interface{})
		values, err := rows.Values()
		if err != nil {
			rows.Close()
			return nil, err
		}
		for i, v := range values {
			v = normalizeValue(v)
			record[t.Fields[i].Name] = v
		}

		rows.Close()

		for _, preload := range toPreload {
			ex, ok := t.Field(preload.Field)
			if !ok {
				return nil, fmt.Errorf("%s is not a field of %s", preload.Field, t.Name)
			}

			var res interface{}
			if _, ok := ex.Annotation("hasOne"); ok {
				res, err = findById(ctx, conn, ex.Type.Type, record[preload.Field], preload.Preload)
				if err != nil {
					return nil, err
				}
			} else if hasMany, ok := ex.Annotation("hasMany"); ok {
				if key, ok := hasMany.Argument("key"); ok {
					keyName := keyField(t)
					res, err = join(ctx, conn, ex.Type.ItemType.Type,
						key.ValueString()+" = $1", []interface{}{record[keyName]},
						preload.Preload)
					if err != nil {
						return nil, err
					}
				}
			}

			record[preload.Field] = res
		}

		return record, nil
	}

	return nil, nil
}

func join(ctx context.Context, conn *pgxpool.Conn, t *spec.Type, where string, args []interface{}, toPreload []Preload) ([]map[string]interface{}, error) {
	sql := generateTableSQL(t)
	if len(where) > 0 {
		sql += " WHERE "
		sql += where
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0, 1000)

	for rows.Next() {
		record := make(map[string]interface{})
		values, err := rows.Values()
		if err != nil {
			rows.Close()
			return nil, err
		}
		for i, v := range values {
			v = normalizeValue(v)
			record[t.Fields[i].Name] = v
		}

		results = append(results, record)
	}

	rows.Close()

	if len(toPreload) > 0 {
		for _, record := range results {
			for _, preload := range toPreload {
				ex, ok := t.Field(preload.Field)
				if !ok {
					return nil, fmt.Errorf("%s is not a field of %s", preload.Field, t.Name)
				}

				var res interface{}
				if _, ok := ex.Annotation("hasOne"); ok {
					res, err = findById(ctx, conn, ex.Type.Type, record[preload.Field], preload.Preload)
					if err != nil {
						return nil, err
					}
				} else if hasMany, ok := ex.Annotation("hasMany"); ok {
					if key, ok := hasMany.Argument("key"); ok {
						keyName := keyField(t)
						res, err = join(ctx, conn, ex.Type.ItemType.Type,
							key.ValueString()+" = $1", []interface{}{record[keyName]},
							preload.Preload)
						if err != nil {
							return nil, err
						}
					}
				}

				record[preload.Field] = res
			}
		}
	}

	return results, nil
}

func getMany(ctx context.Context, conn *pgxpool.Conn, t *spec.Type, input map[string]interface{}, where []Where, toPreload []Preload, offset, limit int64) ([]map[string]interface{}, error) {
	sql := generateTableSQL(t)
	var args []interface{}
	if len(where) > 0 {
		dollarIndex := 1
		for i, part := range where {
			val, err := part.Value.Eval(input)
			if err != nil {
				return nil, err
			}
			if isNil(val) {
				continue
			}
			if i > 0 {
				sql += " AND "
			} else {
				sql += " WHERE "
			}
			query := part.Query
			for strings.Contains(query, "?") {
				query = strings.Replace(query, "?", fmt.Sprintf("$%d", dollarIndex), 1)
				dollarIndex++
			}
			sql += query
			args = append(args, val)
		}
	}
	if offset > 0 {
		sql += " OFFSET " + strconv.FormatInt(offset, 10)
	}
	if limit > 0 {
		sql += " LIMIT " + strconv.FormatInt(limit, 10)
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0, 1000)

	for rows.Next() {
		record := make(map[string]interface{})
		values, err := rows.Values()
		if err != nil {
			rows.Close()
			return nil, err
		}
		for i, v := range values {
			v = normalizeValue(v)
			record[t.Fields[i].Name] = v
		}

		results = append(results, record)
	}

	rows.Close()

	if len(toPreload) > 0 {
		for _, record := range results {
			for _, preload := range toPreload {
				ex, ok := t.Field(preload.Field)
				if !ok {
					return nil, fmt.Errorf("%s is not a field of %s", preload.Field, t.Name)
				}

				var res interface{}
				if _, ok := ex.Annotation("hasOne"); ok {
					res, err = findById(ctx, conn, ex.Type.Type, record[preload.Field], preload.Preload)
					if err != nil {
						return nil, err
					}
				} else if hasMany, ok := ex.Annotation("hasMany"); ok {
					if key, ok := hasMany.Argument("key"); ok {
						keyName := keyField(t)
						res, err = join(ctx, conn, ex.Type.ItemType.Type,
							key.ValueString()+" = $1", []interface{}{record[keyName]},
							preload.Preload)
						if err != nil {
							return nil, err
						}
					}
				}

				record[preload.Field] = res
			}
		}
	}

	return results, nil
}

func getCount(ctx context.Context, conn *pgxpool.Conn, t *spec.Type, input map[string]interface{}, where []Where) (int64, error) {
	sql := generateCountSQL(t)
	var args []interface{}
	if len(where) > 0 {
		dollarIndex := 1
		for i, part := range where {
			val, err := part.Value.Eval(input)
			if err != nil {
				return 0, err
			}
			if isNil(val) {
				continue
			}
			if i > 0 {
				sql += " AND "
			} else {
				sql += " WHERE "
			}
			query := part.Query
			for strings.Contains(query, "?") {
				query = strings.Replace(query, "?", fmt.Sprintf("$%d", dollarIndex), 1)
				dollarIndex++
			}
			sql += query
			args = append(args, val)
		}
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		err = rows.Scan(&count)
	}
	return count, err
}

func keyField(t *spec.Type) string {
	for _, f := range t.Fields {
		if _, ok := f.Annotation("key"); ok {
			return f.Name
		}
	}
	return ""
}

func keyColumn(t *spec.Type) string {
	if _, ok := t.Annotation("primaryKey"); ok {
		return annotationValue(t, "primaryKey", "name", "")
	}
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
			column = annotationValue(f, "hasOne", "foreignKey", "")
		}
		if column == "" {
			if _, ok := f.Annotation("hasMany"); ok {
				column = "1" // Temp solution
			}
		}
		if column == "" {
			column = strcase.ToSnake(f.Name)
		}
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(column)
	}
	buf.WriteString(" FROM ")
	table := annotationValue(t, "entity", "table", t.Name)
	buf.WriteByte('"')
	buf.WriteString(table)
	buf.WriteByte('"')

	return buf.String()
}

func generateCountSQL(t *spec.Type) string {
	var buf strings.Builder

	buf.WriteString("SELECT count(1) FROM ")
	table := annotationValue(t, "entity", "table", t.Name)
	buf.WriteByte('"')
	buf.WriteString(table)
	buf.WriteByte('"')

	return buf.String()
}

func isNil(val interface{}) bool {
	return val == nil ||
		(reflect.ValueOf(val).Kind() == reflect.Ptr &&
			reflect.ValueOf(val).IsNil())
}

func normalizeValue(v interface{}) interface{} {
	switch vv := v.(type) {
	case big.Float:
		v, _ = vv.Float64()
	case big.Int:
		v = vv.Int64()
	case pgtype.Numeric:
		var f float64
		vv.AssignTo(&f)
		v = f
	case pgtype.UUID:
		v = uuid.UUID(vv.Bytes).String()
	case [16]uint8: // UUID
		v = uuid.UUID(vv).String()
	}
	return v
}
