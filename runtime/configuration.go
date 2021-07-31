package runtime

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Specs         []Spec            `json:"specs" yaml:"specs"`
	Compute       Compute           `json:"compute" yaml:"compute"`
	Resiliency    Resiliency        `json:"resiliency" yaml:"resiliency"`
	Services      Services          `json:"services" yaml:"services"`
	Outbound      Services          `json:"outbound" yaml:"outbound"`
	Inbound       FunctionPipelines `json:"inbound" yaml:"inbound"`
	Subscriptions interface{}       `json:"subscriptions" yaml:"subscriptions"`
}

type Spec struct {
	Type string      `json:"type" yaml:"type"`
	With interface{} `json:"with" yaml:"with"`
}

type Compute struct {
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
	Summary        string      `json:"summary" yaml:"summary"`
	Name           string      `json:"name" yaml:"name"`
	With           interface{} `json:"with" yaml:"with"`
	Timeout        string      `json:"timeout" yaml:"timeout"`
	Retry          string      `json:"retry" yaml:"retry"`
	CircuitBreaker string      `json:"circuitBreaker" yaml:"circuitBreaker"`
}

func LoadYAML(in io.Reader) (*Configuration, error) {
	var c Configuration
	if err := yaml.NewDecoder(in).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
