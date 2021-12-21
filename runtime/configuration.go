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
	Specs         []Component                `json:"specs" yaml:"specs"`
	Filters       map[string][]Component     `json:"filters" yaml:"filters"`
	Codecs        map[string]Component       `json:"codecs" yaml:"codecs"`
	Resources     map[string]Component       `json:"resources" yaml:"resources"`
	Compute       Component                  `json:"compute" yaml:"compute"`
	Resiliency    Resiliency                 `json:"resiliency" yaml:"resiliency"`
	Services      Services                   `json:"services" yaml:"services"`
	Providers     Services                   `json:"providers" yaml:"providers"`
	Events        FunctionPipelines          `json:"events" yaml:"events"`
	Flows         FunctionPipelines          `json:"flows" yaml:"flows"`
	Subscriptions interface{}                `json:"subscriptions" yaml:"subscriptions"`
	InputBindings interface{}                `json:"inputBindings" yaml:"inputBindings"`
	Decoding      interface{}                `json:"decoding" yaml:"decoding"`
	Errors        map[string]errorz.Template `json:"errors" yaml:"errors"`
}

type Component struct {
	Type string      `json:"type" yaml:"type"`
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
	Call  string `json:"call" yaml:"call" mapstructure:"call"`
	Steps []Step `json:"steps" yaml:"steps"`
}

type Step struct {
	Name           string      `json:"name" yaml:"name" mapstructure:"name"`
	Call           string      `json:"call" yaml:"call" mapstructure:"call"`
	Uses           string      `json:"uses" yaml:"uses" mapstructure:"uses"`
	With           interface{} `json:"with" yaml:"with" mapstructure:"with"`
	Returns        string      `json:"returns" yaml:"returns" mapstructure:"returns"`
	Timeout        string      `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Retry          string      `json:"retry" yaml:"retry" mapstructure:"retry"`
	CircuitBreaker string      `json:"circuitBreaker" yaml:"circuitBreaker" mapstructure:"circuitBreaker"`
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

		// Flows
		if len(c.Flows) > 0 && config.Flows == nil {
			config.Flows = make(FunctionPipelines)
		}
		for k, v := range c.Flows {
			if _, exists := config.Flows[k]; !exists {
				config.Flows[k] = v
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
