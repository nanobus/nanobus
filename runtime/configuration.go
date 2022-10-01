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

package runtime

import (
	"encoding/json"
	"io"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/nanobus/nanobus/errorz"
)

type Configuration struct {
	Import        []string                   `json:"import" yaml:"import"`
	Transports    map[string]Component       `json:"transports" yaml:"transports"`
	Specs         []Component                `json:"specs" yaml:"specs"`
	Filters       map[string][]Component     `json:"filters" yaml:"filters"`
	Codecs        map[string]Component       `json:"codecs" yaml:"codecs"`
	Resources     map[string]Component       `json:"resources" yaml:"resources"`
	Compute       []Component                `json:"compute" yaml:"compute"`
	Resiliency    Resiliency                 `json:"resiliency" yaml:"resiliency"`
	Services      Services                   `json:"services" yaml:"services"`
	Providers     Services                   `json:"providers" yaml:"providers"`
	Events        FunctionPipelines          `json:"events" yaml:"events"`
	Pipelines     FunctionPipelines          `json:"pipelines" yaml:"pipelines"`
	Subscriptions []Subscription             `json:"subscriptions" yaml:"subscriptions"`
	Errors        map[string]errorz.Template `json:"errors" yaml:"errors"`
}

type Subscription struct {
	Resource  string            `mapstructure:"resource"`
	Topic     string            `mapstructure:"topic"`
	Metadata  map[string]string `mapstructure:"metadata"`
	Codec     string            `mapstructure:"codec"`
	CodecArgs []interface{}     `mapstructure:"codecArgs"`
	Function  string            `mapstructure:"function"`
}

type Component struct {
	Uses string      `json:"uses" yaml:"uses"`
	With interface{} `json:"with" yaml:"with"`
}

type Resiliency struct {
	Timeouts        map[string]Duration         `json:"timeouts" yaml:"timeouts"`
	Retries         map[string]ConfigProperties `json:"retries" yaml:"retries"`
	CircuitBreakers map[string]ConfigProperties `json:"circuitBreakers" yaml:"circuitBreakers"`
}

type ConfigProperties map[string]interface{}

type Services map[string]FunctionPipelines
type FunctionPipelines map[string]Pipeline

type Pipeline struct {
	Name  string `json:"name" yaml:"name"`
	Call  string `json:"call,omitempty" yaml:"call,omitempty" mapstructure:"call"`
	Steps []Step `json:"steps,omitempty" yaml:"steps,omitempty"`
}

type Step struct {
	Name           string      `json:"name" yaml:"name" mapstructure:"name"`
	Call           string      `json:"call,omitempty" yaml:"call,omitempty" mapstructure:"call"`
	Uses           string      `json:"uses,omitempty" yaml:"uses,omitempty" mapstructure:"uses"`
	With           interface{} `json:"with,omitempty" yaml:"with,omitempty" mapstructure:"with"`
	Returns        string      `json:"returns,omitempty" yaml:"returns,omitempty" mapstructure:"returns"`
	Timeout        string      `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout"`
	Retry          string      `json:"retry,omitempty" yaml:"retry,omitempty" mapstructure:"retry"`
	CircuitBreaker string      `json:"circuitBreaker,omitempty" yaml:"circuitBreaker,omitempty" mapstructure:"circuitBreaker"`
	OnError        *Pipeline   `json:"onError,omitempty" yaml:"onError,omitempty" mapstructure:"onError"`
}

func LoadYAML(in io.Reader) (*Configuration, error) {
	var c Configuration
	if err := yaml.NewDecoder(in).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func Combine(config *Configuration, configs ...*Configuration) {
	for _, c := range configs {
		// Filters
		if len(c.Filters) > 0 && config.Filters == nil {
			config.Filters = make(map[string][]Component)
		}
		for k, v := range c.Filters {
			if _, exists := config.Filters[k]; !exists {
				config.Filters[k] = v
			}
		}

		// Codecs
		if len(c.Codecs) > 0 && config.Codecs == nil {
			config.Codecs = make(map[string]Component)
		}
		for k, v := range c.Codecs {
			if _, exists := config.Codecs[k]; !exists {
				config.Codecs[k] = v
			}
		}

		// Resiliency
		if len(c.Resiliency.Timeouts) > 0 && config.Resiliency.Timeouts == nil {
			config.Resiliency.Timeouts = make(map[string]Duration)
		}
		for k, v := range c.Resiliency.Timeouts {
			if _, exists := config.Resiliency.Timeouts[k]; !exists {
				config.Resiliency.Timeouts[k] = v
			}
		}

		if len(c.Resiliency.Retries) > 0 && config.Resiliency.Retries == nil {
			config.Resiliency.Retries = make(map[string]ConfigProperties)
		}
		for k, v := range c.Resiliency.Retries {
			if _, exists := config.Resiliency.Retries[k]; !exists {
				config.Resiliency.Retries[k] = v
			}
		}

		if len(c.Resiliency.CircuitBreakers) > 0 && config.Resiliency.CircuitBreakers == nil {
			config.Resiliency.CircuitBreakers = make(map[string]ConfigProperties)
		}
		for k, v := range c.Resiliency.CircuitBreakers {
			if _, exists := config.Resiliency.CircuitBreakers[k]; !exists {
				config.Resiliency.CircuitBreakers[k] = v
			}
		}

		// Services
		if len(c.Services) > 0 && config.Services == nil {
			config.Services = make(Services)
		}
		for k, v := range c.Services {
			existing, exists := config.Services[k]
			if !exists {
				existing = make(FunctionPipelines)
				config.Services[k] = existing
			}
			for k, v := range v {
				if _, exists := existing[k]; !exists {
					existing[k] = v
				}
			}
		}

		// Providers
		if len(c.Providers) > 0 && config.Providers == nil {
			config.Providers = make(Services)
		}
		for k, v := range c.Providers {
			existing, exists := config.Providers[k]
			if !exists {
				existing = make(FunctionPipelines)
				config.Providers[k] = existing
			}
			for k, v := range v {
				if _, exists := existing[k]; !exists {
					existing[k] = v
				}
			}
		}

		// Events
		if len(c.Events) > 0 && config.Events == nil {
			config.Events = make(FunctionPipelines)
		}
		for k, v := range c.Events {
			if _, exists := config.Events[k]; !exists {
				config.Events[k] = v
			}
		}

		// Pipelines
		if len(c.Pipelines) > 0 && config.Pipelines == nil {
			config.Pipelines = make(FunctionPipelines)
		}
		for k, v := range c.Pipelines {
			if _, exists := config.Pipelines[k]; !exists {
				config.Pipelines[k] = v
			}
		}

		// Errors
		if len(c.Errors) > 0 && config.Errors == nil {
			config.Errors = make(map[string]errorz.Template)
		}
		for k, v := range c.Errors {
			if _, exists := config.Errors[k]; !exists {
				config.Errors[k] = v
			}
		}
	}
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	return d.Parse(str)
}

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	return d.Parse(str)
}

func (d *Duration) Parse(str string) error {
	millis, err := strconv.ParseUint(str, 10, 32)
	if err == nil {
		*d = Duration(millis) * Duration(time.Millisecond)
		return nil
	}

	dur, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	*d = Duration(dur)

	return nil
}
