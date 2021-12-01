package errorz

import (
	"encoding/json"
	"strings"
	"text/template"
)

type Resolver func(err error) *Error

type Template struct {
	Type      string             `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type"`
	Code      ErrCode            `json:"code" yaml:"code" mapstructure:"code"`
	Status    int                `json:"status,omitempty" yaml:"status,omitempty" mapstructure:"status"`
	Title     *Text              `json:"title,omitempty" yaml:"title,omitempty" mapstructure:"title"`
	Message   *Text              `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message"`
	Instance  *Text              `json:"instance,omitempty" yaml:"instance,omitempty" mapstructure:"instance"`
	Languages map[string]Strings `json:"languages,omitempty" yaml:"languages,omitempty" mapstructure:"languages"`
}

type Strings struct {
	Title   *Text `json:"title,omitempty" yaml:"title,omitempty" mapstructure:"title"`
	Message *Text `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message"`
}

type Text struct {
	tmpl *template.Template
}

func (t *Text) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	return t.Parse(str)
}

func (t *Text) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	return t.Parse(str)
}

func (t *Text) Parse(str string) error {
	tmpl, err := template.New("text").Parse(str)
	if err != nil {
		return err
	}

	*t = Text{
		tmpl: tmpl,
	}

	return nil
}

func (t *Text) Eval(data interface{}) (string, error) {
	var out strings.Builder
	if err := t.tmpl.Execute(&out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}
