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

package expr

import (
	"encoding/json"
	"reflect"
	"strings"
	"text/template"
)

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
	tmpl, err := template.New("text").Funcs(template.FuncMap{
		"pick": pick,
	}).Parse(str)
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

func pick(args ...interface{}) interface{} {
	for _, arg := range args {
		if !isNil(arg) {
			return arg
		}
	}
	return ""
}

func isNil(val interface{}) bool {
	return val == nil ||
		(reflect.ValueOf(val).Kind() == reflect.Ptr &&
			reflect.ValueOf(val).IsNil())
}
