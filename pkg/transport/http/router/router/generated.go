// Code generated by @apexlang/codegen. DO NOT EDIT.

package router

import (
	"github.com/nanobus/nanobus/pkg/handler"
	"github.com/nanobus/nanobus/pkg/transport/http/router"
)

type RouterV1Config []Route

func RouterV1() (string, router.Loader) {
	return "nanobus.transport.http.router/v1", RouterV1Loader
}

type Route struct {
	Methods  string          `json:"methods" yaml:"methods" msgpack:"methods" mapstructure:"methods" validate:"required"`
	URI      string          `json:"uri" yaml:"uri" msgpack:"uri" mapstructure:"uri" validate:"required"`
	Encoding *string         `json:"encoding,omitempty" yaml:"encoding,omitempty" msgpack:"encoding,omitempty" mapstructure:"encoding"`
	Handler  handler.Handler `json:"handler" yaml:"handler" msgpack:"handler" mapstructure:"handler" validate:"required"`
}