// Code generated by @apexlang/codegen. DO NOT EDIT.

package time

import (
	"github.com/nanobus/nanobus/pkg/handler"
	"github.com/nanobus/nanobus/pkg/transport"
)

type TimeSchedulerV1Config struct {
	Schedules []Schedule `json:"schedules" yaml:"schedules" msgpack:"schedules" mapstructure:"schedules" validate:"dive"`
}

func TimeSchedulerV1() (string, transport.Loader) {
	return "nanobus.transport.time.scheduler/v1", TimeSchedulerV1Loader
}

type Schedule struct {
	Schedule string          `json:"schedule" yaml:"schedule" msgpack:"schedule" mapstructure:"schedule" validate:"required"`
	Handler  handler.Handler `json:"handler" yaml:"handler" msgpack:"handler" mapstructure:"handler" validate:"required"`
}
