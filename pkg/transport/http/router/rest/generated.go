// Code generated by @apexlang/codegen. DO NOT EDIT.

package rest

import (
	"github.com/nanobus/nanobus/pkg/transport/http/router"
)

type RestV1Config struct {
	Documentation Documentation `json:"documentation" yaml:"documentation" msgpack:"documentation" mapstructure:"documentation" validate:"required"`
}

func RestV1() (string, router.Loader) {
	return "nanobus.transport.http.rest/v1", RestV1Loader
}

type Documentation struct {
	SwaggerUI  bool `json:"swaggerUI" yaml:"swaggerUI" msgpack:"swaggerUI" mapstructure:"swaggerUI" validate:"required"`
	Postman    bool `json:"postman" yaml:"postman" msgpack:"postman" mapstructure:"postman" validate:"required"`
	RestClient bool `json:"restClient" yaml:"restClient" msgpack:"restClient" mapstructure:"restClient" validate:"required"`
}