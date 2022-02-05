package postgres

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"strings"

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

func getOne(ctx context.Context, conn *pgxpool.Conn, t *spec.Type, idValue interface{}, toPreload []Preload) (interface{}, error) {
	idColumn := keyColumn(t)
	sql := generateTableSQL(t) + " WHERE " + idColumn + "=$1"
	fmt.Println(sql, idValue)
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
			record[t.Fields[i].Name] = v
		}

		rows.Close()

		for _, preload := range toPreload {
			ex, ok := t.Field(preload.Field)
			if !ok {
				return nil, fmt.Errorf("%s is not a field of %s", preload.Field, t.Name)
			}

			res, err := getOne(ctx, conn, ex.Type.Type, record[preload.Field], preload.Preload)
			if err != nil {
				return nil, err
			}

			record[preload.Field] = res
		}

		return record, nil
	}

	return nil, nil
}

func getMany(ctx context.Context, conn *pgxpool.Conn, t *spec.Type, input map[string]interface{}, where []Where, toPreload []Preload) ([]map[string]interface{}, error) {
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
	fmt.Println(sql)
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0, 1000)

	if rows.Next() {
		record := make(map[string]interface{})
		values, err := rows.Values()
		if err != nil {
			rows.Close()
			return nil, err
		}
		for i, v := range values {
			switch vv := v.(type) {
			case big.Float:
				v, _ = vv.Float64()
			case big.Int:
				v = vv.Int64()
			case pgtype.Numeric:
				var f float64
				vv.AssignTo(&f)
				v = f
			}
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

				res, err := getOne(ctx, conn, ex.Type.Type, record[preload.Field], preload.Preload)
				if err != nil {
					return nil, err
				}

				record[preload.Field] = res
			}
		}
	}

	return results, nil
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
			column = annotationValue(f, "hasOne", "foreignKey", "")
			if column == "" {
				continue
			}
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

func isNil(val interface{}) bool {
	return val == nil ||
		(reflect.ValueOf(val).Kind() == reflect.Ptr &&
			reflect.ValueOf(val).IsNil())
}
