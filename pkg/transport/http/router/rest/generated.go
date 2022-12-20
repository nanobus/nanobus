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
	SwaggerUI  *bool `json:"swaggerUI,omitempty" yaml:"swaggerUI,omitempty" msgpack:"swaggerUI,omitempty" mapstructure:"swaggerUI"`
	Postman    *bool `json:"postman,omitempty" yaml:"postman,omitempty" msgpack:"postman,omitempty" mapstructure:"postman"`
	RestClient *bool `json:"restClient,omitempty" yaml:"restClient,omitempty" msgpack:"restClient,omitempty" mapstructure:"restClient"`
}
