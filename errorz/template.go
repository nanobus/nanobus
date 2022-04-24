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

package errorz

import (
	"github.com/nanobus/nanobus/expr"
)

type Resolver func(err error) *Error

type Template struct {
	Type    string             `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type"`
	Code    ErrCode            `json:"code" yaml:"code" mapstructure:"code"`
	Status  int                `json:"status,omitempty" yaml:"status,omitempty" mapstructure:"status"`
	Title   *expr.Text         `json:"title,omitempty" yaml:"title,omitempty" mapstructure:"title"`
	Message *expr.Text         `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message"`
	Path    string             `json:"path,omitempty" yaml:"path,omitempty" mapstructure:"path"`
	Help    *expr.Text         `json:"help,omitempty" yaml:"help,omitempty" mapstructure:"help"`
	Locales map[string]Strings `json:"locales,omitempty" yaml:"locales,omitempty" mapstructure:"locales"`
}

type Strings struct {
	Title   *expr.Text `json:"title,omitempty" yaml:"title,omitempty" mapstructure:"title"`
	Message *expr.Text `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message"`
}
