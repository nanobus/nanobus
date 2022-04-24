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

package resolve

import (
	"errors"
	"fmt"
	"reflect"
)

type DependencyResolver func(name string) (interface{}, bool)

type ResolveAs func(name string, target interface{}) bool

func Resolve(resolver ResolveAs, args ...interface{}) error {
	if len(args)%2 != 0 {
		return errors.New("invalid number of arguments passed to Resolve")
	}

	for i := 0; i < len(args); i += 2 {
		dependencyName, ok := args[i].(string)
		if !ok {
			return fmt.Errorf("argument %d is not a string", i)
		}

		if !resolver(dependencyName, args[i+1]) {
			return fmt.Errorf("could not resolve dependency %q", dependencyName)
		}
	}

	return nil
}

func ToResolveAs(resolver DependencyResolver) ResolveAs {
	return func(name string, target interface{}) bool {
		dependency, ok := resolver(name)
		if !ok {
			return false
		}

		return As(dependency, target)
	}
}

func As(source, target interface{}) bool {
	if target == nil {
		return false
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}

	targetType := typ.Elem()
	if reflect.TypeOf(source).AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(source))
		return true
	}

	return false
}
