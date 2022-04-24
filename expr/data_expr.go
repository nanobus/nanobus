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
	"errors"
	"fmt"

	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
)

type DataExpr struct {
	script string
}

var ErrNotAMap = errors.New("data expression result was not a map")

func (d *DataExpr) DecodeString(value string) error {
	d.script = value

	return nil
}

func (d *DataExpr) Expr() string {
	return d.script
}

func (d *DataExpr) Eval(variables map[string]interface{}) (interface{}, error) {
	result, err := d.doEval(variables)
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[interface{}]interface{}); ok {
		result = mIItoSI(m)
	}

	return result, err
}

func (d *DataExpr) EvalMap(variables map[string]interface{}) (map[string]string, error) {
	result, err := d.doEval(variables)
	if err != nil {
		return nil, err
	}

	switch t := result.(type) {
	case map[interface{}]interface{}:
		return mIItoSS(t), nil
	case map[string]interface{}:
		return mSItoSS(t), nil
	case map[string]string:
		return t, nil
	}

	return nil, ErrNotAMap
}

func (d *DataExpr) doEval(variables map[string]interface{}) (interface{}, error) {
	e := env.NewEnv()

	err := e.Define("println", fmt.Println)
	if err != nil {
		return nil, fmt.Errorf("define error: %w", err)
	}
	for name, value := range variables {
		err = e.Define(name, value)
		if err != nil {
			return nil, fmt.Errorf("define error: %w", err)
		}
	}

	result, err := vm.Execute(e, nil, d.script)
	if err != nil {
		return nil, fmt.Errorf("execute error: %w", err)
	}

	return result, nil
}

func mIItoSI(m map[interface{}]interface{}) map[string]interface{} {
	ret := make(map[string]interface{}, len(m))
	for k, v := range m {
		if km, ok := v.(map[interface{}]interface{}); ok {
			v = mIItoSI(km)
		}
		ret[interfaceToString(k)] = v
	}
	return ret
}

func mIItoSS(m map[interface{}]interface{}) map[string]string {
	ret := make(map[string]string, len(m))
	for k, v := range m {
		switch km := v.(type) {
		case map[interface{}]interface{}:
			v = mIItoSS(km)
		case map[string]interface{}:
			v = mSItoSS(km)
		}
		ret[interfaceToString(k)] = interfaceToString(v)
	}
	return ret
}

func mSItoSS(m map[string]interface{}) map[string]string {
	ret := make(map[string]string, len(m))
	for k, v := range m {
		switch km := v.(type) {
		case map[interface{}]interface{}:
			v = mIItoSS(km)
		case map[string]interface{}:
			v = mSItoSS(km)
		}
		ret[k] = interfaceToString(v)
	}
	return ret
}

func interfaceToString(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	if s, ok := value.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", value)
}
