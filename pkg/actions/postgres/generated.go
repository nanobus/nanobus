// Code generated by @apexlang/codegen. DO NOT EDIT.

package postgres

import (
	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/expr"
)

type ExecConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Data is the input bindings sent.
	Data *expr.DataExpr `json:"Data,omitempty" yaml:"Data,omitempty" msgpack:"Data,omitempty" mapstructure:"Data"`
	// SQL is the SQL query to execute.
	SQL string `json:"sql" yaml:"sql" msgpack:"sql" mapstructure:"sql" validate:"required"`
	// Args are the evaluations to use as arguments for the SQL query.
	Args []*expr.ValueExpr `json:"args,omitempty" yaml:"args,omitempty" msgpack:"args,omitempty" mapstructure:"args"`
}

func Exec() (string, actions.Loader) {
	return "@postgres/exec", ExecLoader
}

type ExecMultiConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Statements are the statements to execute within a single transaction.
	Statements []Statement `json:"statements" yaml:"statements" msgpack:"statements" mapstructure:"statements" validate:"required"`
}

func ExecMulti() (string, actions.Loader) {
	return "@postgres/exec_multi", ExecMultiLoader
}

type Statement struct {
	// Data is the input bindings sent.
	Data *expr.DataExpr `json:"data,omitempty" yaml:"data,omitempty" msgpack:"data,omitempty" mapstructure:"data"`
	// SQL is the SQL query to execute.
	SQL string `json:"sql" yaml:"sql" msgpack:"sql" mapstructure:"sql" validate:"required"`
	// Args are the evaluations to use as arguments for the SQL query.
	Args []*expr.ValueExpr `json:"args,omitempty" yaml:"args,omitempty" msgpack:"args,omitempty" mapstructure:"args"`
}

type FindOneConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Namespace is the type namespace to load.
	Namespace string `json:"namespace" yaml:"namespace" msgpack:"namespace" mapstructure:"namespace" validate:"required"`
	// Type is the type name to load.
	Type string `json:"type" yaml:"type" msgpack:"type" mapstructure:"type" validate:"required"`
	// Preload lists the relationship to expand/load.
	Preload []Preload `json:"preload,omitempty" yaml:"preload,omitempty" msgpack:"preload,omitempty" mapstructure:"preload"`
	// Where list the parts of the where clause.
	Where []Where `json:"where,omitempty" yaml:"where,omitempty" msgpack:"where,omitempty" mapstructure:"where"`
	// NotFoundError is the error to return if the key is not found.
	NotFoundError string `json:"notFoundError" yaml:"notFoundError" msgpack:"notFoundError" mapstructure:"notFoundError" validate:"required"`
}

func FindOne() (string, actions.Loader) {
	return "@postgres/find_one", FindOneLoader
}

type Preload struct {
	Field   string    `json:"field" yaml:"field" msgpack:"field" mapstructure:"field" validate:"required"`
	Preload []Preload `json:"preload" yaml:"preload" msgpack:"preload" mapstructure:"preload" validate:"required"`
}

type Where struct {
	Query string          `json:"query" yaml:"query" msgpack:"query" mapstructure:"query" validate:"required"`
	Value *expr.ValueExpr `json:"value" yaml:"value" msgpack:"value" mapstructure:"value" validate:"required"`
}

type FindConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Namespace is the type namespace to load.
	Namespace string `json:"namespace" yaml:"namespace" msgpack:"namespace" mapstructure:"namespace" validate:"required"`
	// Type is the type name to load.
	Type string `json:"type" yaml:"type" msgpack:"type" mapstructure:"type" validate:"required"`
	// Preload lists the relationship to expand/load.
	Preload []Preload `json:"preload,omitempty" yaml:"preload,omitempty" msgpack:"preload,omitempty" mapstructure:"preload"`
	// Where list the parts of the where clause.
	Where []Where `json:"where,omitempty" yaml:"where,omitempty" msgpack:"where,omitempty" mapstructure:"where"`
	// Pagination is the optional fields to wrap the results with.
	Pagination *Pagination `json:"pagination,omitempty" yaml:"pagination,omitempty" msgpack:"pagination,omitempty" mapstructure:"pagination"`
	// Offset is the query offset.
	Offset *expr.ValueExpr `json:"offset,omitempty" yaml:"offset,omitempty" msgpack:"offset,omitempty" mapstructure:"offset"`
	// Limit is the query limit.
	Limit *expr.ValueExpr `json:"limit,omitempty" yaml:"limit,omitempty" msgpack:"limit,omitempty" mapstructure:"limit"`
}

func Find() (string, actions.Loader) {
	return "@postgres/find", FindLoader
}

type Pagination struct {
	PageIndex *string `json:"pageIndex,omitempty" yaml:"pageIndex,omitempty" msgpack:"pageIndex,omitempty" mapstructure:"pageIndex"`
	PageCount *string `json:"pageCount,omitempty" yaml:"pageCount,omitempty" msgpack:"pageCount,omitempty" mapstructure:"pageCount"`
	Offset    *string `json:"offset,omitempty" yaml:"offset,omitempty" msgpack:"offset,omitempty" mapstructure:"offset"`
	Limit     *string `json:"limit,omitempty" yaml:"limit,omitempty" msgpack:"limit,omitempty" mapstructure:"limit"`
	Count     *string `json:"count,omitempty" yaml:"count,omitempty" msgpack:"count,omitempty" mapstructure:"count"`
	Total     *string `json:"total,omitempty" yaml:"total,omitempty" msgpack:"total,omitempty" mapstructure:"total"`
	Items     string  `json:"items" yaml:"items" msgpack:"items" mapstructure:"items" validate:"required"`
}

type LoadConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Namespace is the type namespace to load.
	Namespace string `json:"namespace" yaml:"namespace" msgpack:"namespace" mapstructure:"namespace" validate:"required"`
	// Type is the type name to load.
	Type string `json:"type" yaml:"type" msgpack:"type" mapstructure:"type" validate:"required"`
	// ID is the entity identifier expression.
	Key *expr.ValueExpr `json:"key" yaml:"key" msgpack:"key" mapstructure:"key" validate:"required"`
	// Preload lists the relationship to expand/load.
	Preload []Preload `json:"preload,omitempty" yaml:"preload,omitempty" msgpack:"preload,omitempty" mapstructure:"preload"`
	// NotFoundError is the error to return if the key is not found.
	NotFoundError string `json:"notFoundError" yaml:"notFoundError" msgpack:"notFoundError" mapstructure:"notFoundError" validate:"required"`
}

func Load() (string, actions.Loader) {
	return "@postgres/load", LoadLoader
}

type QueryConfig struct {
	// Resource is the name of the connection resource to use.
	Resource string `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// SQL is the SQL query to execute.
	SQL string `json:"sql" yaml:"sql" msgpack:"sql" mapstructure:"sql" validate:"required"`
	// Args are the evaluations to use as arguments for the SQL query.
	Args []*expr.ValueExpr `json:"args,omitempty" yaml:"args,omitempty" msgpack:"args,omitempty" mapstructure:"args"`
	// Single indicates a single row should be returned if found.
	Single bool `json:"single" yaml:"single" msgpack:"single" mapstructure:"single" validate:"required"`
}

func Query() (string, actions.Loader) {
	return "@postgres/query", QueryLoader
}
