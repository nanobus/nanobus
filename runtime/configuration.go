package runtime

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Specs         []Component            `json:"specs" yaml:"specs"`
	Filters       map[string][]Component `json:"filters" yaml:"filters"`
	Codecs        map[string]Component   `json:"codecs" yaml:"codecs"`
	Compute       Component              `json:"compute" yaml:"compute"`
	Resiliency    Resiliency             `json:"resiliency" yaml:"resiliency"`
	Services      Services               `json:"services" yaml:"services"`
	Providers     Services               `json:"providers" yaml:"providers"`
	Events        FunctionPipelines      `json:"events" yaml:"events"`
	Subscriptions interface{}            `json:"subscriptions" yaml:"subscriptions"`
	InputBindings interface{}            `json:"inputBindings" yaml:"inputBindings"`
	Decoding      interface{}            `json:"decoding" yaml:"decoding"`
}

type Component struct {
	Type string      `json:"type" yaml:"type"`
	With interface{} `json:"with" yaml:"with"`
}

type Resiliency struct {
	Retries         map[string]map[string]interface{} `json:"retries" yaml:"retries"`
	CircuitBreakers map[string]map[string]interface{} `json:"circuitBreakers" yaml:"circuitBreakers"`
}

type Services map[string]FunctionPipelines
type FunctionPipelines map[string]Pipeline

type Pipeline struct {
	Summary string `json:"summary" yaml:"summary"`
	Actions []Step `json:"actions" yaml:"actions"`
}

type Step struct {
	Summary        string      `json:"summary" yaml:"summary" mapstructure:"summary"`
	Name           string      `json:"name" yaml:"name" mapstructure:"name"`
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
