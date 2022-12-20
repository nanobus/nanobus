// Code generated by @apexlang/codegen. DO NOT EDIT.

package time

import (
	"github.com/nanobus/nanobus/pkg/runtime"
	"github.com/nanobus/nanobus/pkg/transport"
)

type TimeSchedulerV1Config struct {
	Schedule string              `json:"schedule" yaml:"schedule" msgpack:"schedule" mapstructure:"schedule" validate:"required"`
	Action   []runtime.Component `json:"action,omitempty" yaml:"action,omitempty" msgpack:"action,omitempty" mapstructure:"action"`
}

func TimeSchedulerV1() (string, transport.Loader) {
	return "nanobus.transport.time.scheduler/v1", TimeSchedulerV1Loader
}
