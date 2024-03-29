// Code generated by @apexlang/codegen. DO NOT EDIT.

package dapr

import (
	"encoding/json"
	"errors"

	"github.com/nanobus/nanobus/pkg/actions"
	"github.com/nanobus/nanobus/pkg/expr"
	"github.com/nanobus/nanobus/pkg/handler"
	"github.com/nanobus/nanobus/pkg/resource"
)

type CodecRef string

// TODO
type InvokeBindingConfig struct {
	// The name of the Dapr client resource.
	Resource resource.Ref `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Name of binding to invoke.
	Binding string `json:"binding" yaml:"binding" msgpack:"binding" mapstructure:"binding" validate:"required"`
	// Name of the operation type for the binding to invoke.
	Operation string `json:"operation" yaml:"operation" msgpack:"operation" mapstructure:"operation" validate:"required"`
	// The configured codec to use for encoding the message.
	Codec CodecRef `json:"codec" yaml:"codec" msgpack:"codec" mapstructure:"codec"`
	// The arguments for the codec, if any.
	CodecArgs []interface{} `json:"codecArgs,omitempty" yaml:"codecArgs,omitempty" msgpack:"codecArgs,omitempty" mapstructure:"codecArgs" validate:"dive"`
	// Data is the input data sent.
	Data *expr.DataExpr `json:"data,omitempty" yaml:"data,omitempty" msgpack:"data,omitempty" mapstructure:"data"`
	// Metadata is the input binding metadata.
	Metadata *expr.DataExpr `json:"metadata,omitempty" yaml:"metadata,omitempty" msgpack:"metadata,omitempty" mapstructure:"metadata"`
}

func InvokeBinding() (string, actions.Loader) {
	return "@dapr/invoke_binding", InvokeBindingLoader
}

// TODO
type PublishConfig struct {
	// The name of the Dapr client resource.
	Resource resource.Ref `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Name of pubsub to invoke.
	Pubsub string `json:"pubsub" yaml:"pubsub" msgpack:"pubsub" mapstructure:"pubsub" validate:"required"`
	// Topic is the name of the topic to publish to.
	Topic string `json:"topic" yaml:"topic" msgpack:"topic" mapstructure:"topic" validate:"required"`
	// The configured codec to use for encoding the message.
	Codec CodecRef `json:"codec" yaml:"codec" msgpack:"codec" mapstructure:"codec"`
	// The arguments for the codec, if any.
	CodecArgs []interface{} `json:"codecArgs,omitempty" yaml:"codecArgs,omitempty" msgpack:"codecArgs,omitempty" mapstructure:"codecArgs" validate:"dive"`
	// optional value to use for the message key (is supported).
	Key *expr.ValueExpr `json:"key,omitempty" yaml:"key,omitempty" msgpack:"key,omitempty" mapstructure:"key"`
	// The input bindings sent.
	Data *expr.DataExpr `json:"data,omitempty" yaml:"data,omitempty" msgpack:"data,omitempty" mapstructure:"data"`
	// The input binding metadata.
	Metadata *expr.DataExpr `json:"metadata,omitempty" yaml:"metadata,omitempty" msgpack:"metadata,omitempty" mapstructure:"metadata"`
	// Enables/disables propogating the distributed tracing context (e.g. W3C
	// TraceContext standard).
	PropogateTracing bool `json:"propogateTracing" yaml:"propogateTracing" msgpack:"propogateTracing" mapstructure:"propogateTracing"`
}

func Publish() (string, actions.Loader) {
	return "@dapr/publish", PublishLoader
}

// TODO
type DeleteStateConfig struct {
	// The name of the Dapr client resource.
	Resource resource.Ref `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Name of state store to invoke.
	Store string `json:"store" yaml:"store" msgpack:"store" mapstructure:"store" validate:"required"`
	// The key to delete.
	Key *expr.ValueExpr `json:"key" yaml:"key" msgpack:"key" mapstructure:"key" validate:"required"`
	// Etag value of the item to delete
	Etag *expr.ValueExpr `json:"etag,omitempty" yaml:"etag,omitempty" msgpack:"etag,omitempty" mapstructure:"etag"`
	// The desired concurrency level
	Concurrency Concurrency `json:"concurrency" yaml:"concurrency" msgpack:"concurrency" mapstructure:"concurrency"`
	// The desired consistency level
	Consistency Consistency `json:"consistency" yaml:"consistency" msgpack:"consistency" mapstructure:"consistency"`
}

func DeleteState() (string, actions.Loader) {
	return "@dapr/delete_state", DeleteStateLoader
}

// TODO
type GetStateConfig struct {
	// The name of the Dapr client resource.
	Resource resource.Ref `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Name of state store to invoke.
	Store string `json:"store" yaml:"store" msgpack:"store" mapstructure:"store" validate:"required"`
	// The key to get.
	Key *expr.ValueExpr `json:"key" yaml:"key" msgpack:"key" mapstructure:"key" validate:"required"`
	// The configured codec to use for decoding the state.
	Codec CodecRef `json:"codec" yaml:"codec" msgpack:"codec" mapstructure:"codec"`
	// The arguments for the codec, if any.
	CodecArgs []interface{} `json:"codecArgs,omitempty" yaml:"codecArgs,omitempty" msgpack:"codecArgs,omitempty" mapstructure:"codecArgs" validate:"dive"`
	// The error to return if the key is not found.
	NotFoundError string `json:"notFoundError" yaml:"notFoundError" msgpack:"notFoundError" mapstructure:"notFoundError" validate:"required"`
	// The desired concurrency level
	Concurrency Concurrency `json:"concurrency" yaml:"concurrency" msgpack:"concurrency" mapstructure:"concurrency"`
	// The desired consistency level
	Consistency Consistency `json:"consistency" yaml:"consistency" msgpack:"consistency" mapstructure:"consistency"`
}

func GetState() (string, actions.Loader) {
	return "@dapr/get_state", GetStateLoader
}

// TODO
type SetStateConfig struct {
	// The name of the Dapr client resource.
	Resource resource.Ref `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// Name of state store to invoke.
	Store string `json:"store" yaml:"store" msgpack:"store" mapstructure:"store" validate:"required"`
	// The configured codec to use for encoding the state.
	Codec CodecRef `json:"codec" yaml:"codec" msgpack:"codec" mapstructure:"codec"`
	// The arguments for the codec, if any.
	CodecArgs []interface{} `json:"codecArgs,omitempty" yaml:"codecArgs,omitempty" msgpack:"codecArgs,omitempty" mapstructure:"codecArgs" validate:"dive"`
	// The items to set in the store.
	Items []SetStateItem `json:"items" yaml:"items" msgpack:"items" mapstructure:"items" validate:"dive"`
}

func SetState() (string, actions.Loader) {
	return "@dapr/set_state", SetStateLoader
}

// TODO
type SetStateItem struct {
	// The key of the item to set.
	Key *expr.ValueExpr `json:"key" yaml:"key" msgpack:"key" mapstructure:"key" validate:"required"`
	// an option expression to evaluate a.
	ForEach *expr.ValueExpr `json:"forEach,omitempty" yaml:"forEach,omitempty" msgpack:"forEach,omitempty" mapstructure:"forEach"`
	// Optional data expression to tranform the data to set.
	Value *expr.DataExpr `json:"value,omitempty" yaml:"value,omitempty" msgpack:"value,omitempty" mapstructure:"value"`
	// Etag value of the item to set
	Etag *expr.ValueExpr `json:"etag,omitempty" yaml:"etag,omitempty" msgpack:"etag,omitempty" mapstructure:"etag"`
	// Optional data expression for the key's metadata.
	Metadata *expr.DataExpr `json:"metadata,omitempty" yaml:"metadata,omitempty" msgpack:"metadata,omitempty" mapstructure:"metadata"`
	// The desired concurrency level
	Concurrency Concurrency `json:"concurrency" yaml:"concurrency" msgpack:"concurrency" mapstructure:"concurrency"`
	// The desired consistency level
	Consistency Consistency `json:"consistency" yaml:"consistency" msgpack:"consistency" mapstructure:"consistency"`
}

type InvokeActorConfig struct {
	// The name of the Dapr client resource.
	Resource resource.Ref `json:"resource" yaml:"resource" msgpack:"resource" mapstructure:"resource" validate:"required"`
	// The actor handler (type::method)
	Handler handler.Handler `json:"handler" yaml:"handler" msgpack:"handler" mapstructure:"handler" validate:"required"`
	// The actor identifier
	ID *expr.ValueExpr `json:"id" yaml:"id" msgpack:"id" mapstructure:"id" validate:"required"`
	// The input sent.
	Data *expr.DataExpr `json:"data,omitempty" yaml:"data,omitempty" msgpack:"data,omitempty" mapstructure:"data"`
	// The configured codec to use for encoding the message.
	Codec CodecRef `json:"codec" yaml:"codec" msgpack:"codec" mapstructure:"codec"`
	// The arguments for the codec, if any.
	CodecArgs []interface{} `json:"codecArgs,omitempty" yaml:"codecArgs,omitempty" msgpack:"codecArgs,omitempty" mapstructure:"codecArgs" validate:"dive"`
}

func InvokeActor() (string, actions.Loader) {
	return "@dapr/invoke_actor", InvokeActorLoader
}

// TODO
type Concurrency int32

const (
	// Undefined value for state concurrency
	ConcurrencyUndefined Concurrency = 0
	// First write concurrency value
	ConcurrencyFirstWrite Concurrency = 1
	// Last write concurrency value
	ConcurrencyLastWrite Concurrency = 2
)

var toStringConcurrency = map[Concurrency]string{
	ConcurrencyUndefined:  "undefined",
	ConcurrencyFirstWrite: "firstWrite",
	ConcurrencyLastWrite:  "lastWrite",
}

var toIDConcurrency = map[string]Concurrency{
	"undefined":  ConcurrencyUndefined,
	"firstWrite": ConcurrencyFirstWrite,
	"lastWrite":  ConcurrencyLastWrite,
}

func (e Concurrency) String() string {
	str, ok := toStringConcurrency[e]
	if !ok {
		return "unknown"
	}
	return str
}

func (e *Concurrency) FromString(str string) error {
	var ok bool
	*e, ok = toIDConcurrency[str]
	if !ok {
		return errors.New("unknown value \"" + str + "\" for Concurrency")
	}
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (e Concurrency) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *Concurrency) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	return e.FromString(str)
}

// TODO
type Consistency int32

const (
	// Undefined value for state consistency
	ConsistencyUndefined Consistency = 0
	// Eventual state consistency value
	ConsistencyEventual Consistency = 1
	// Strong state consistency value
	ConsistencyStrong Consistency = 2
)

var toStringConsistency = map[Consistency]string{
	ConsistencyUndefined: "undefined",
	ConsistencyEventual:  "eventual",
	ConsistencyStrong:    "strong",
}

var toIDConsistency = map[string]Consistency{
	"undefined": ConsistencyUndefined,
	"eventual":  ConsistencyEventual,
	"strong":    ConsistencyStrong,
}

func (e Consistency) String() string {
	str, ok := toStringConsistency[e]
	if !ok {
		return "unknown"
	}
	return str
}

func (e *Consistency) FromString(str string) error {
	var ok bool
	*e, ok = toIDConsistency[str]
	if !ok {
		return errors.New("unknown value \"" + str + "\" for Consistency")
	}
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (e Consistency) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *Consistency) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	return e.FromString(str)
}
