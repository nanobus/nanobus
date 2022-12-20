// Code generated by @apexlang/codegen. DO NOT EDIT.

package router

import (
	"github.com/nanobus/nanobus/pkg/handler"
	"github.com/nanobus/nanobus/pkg/transport/http/router"
)

type RouterV1Config struct {
	Routes []AddRoute `json:"routes" yaml:"routes" msgpack:"routes" mapstructure:"routes" validate:"required"`
}

func RouterV1() (string, router.Loader) {
	return "nanobus.transport.http.router/v1", RouterV1Loader
}

type AddRoute struct {
	Method   string          `json:"method" yaml:"method" msgpack:"method" mapstructure:"method" validate:"required"`
	URI      string          `json:"uri" yaml:"uri" msgpack:"uri" mapstructure:"uri" validate:"required"`
	Encoding *string         `json:"encoding,omitempty" yaml:"encoding,omitempty" msgpack:"encoding,omitempty" mapstructure:"encoding"`
	Handler  handler.Handler `json:"handler" yaml:"handler" msgpack:"handler" mapstructure:"handler" validate:"required"`
}
