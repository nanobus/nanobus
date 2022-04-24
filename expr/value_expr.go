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
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	expr_proto "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type ValueExpr struct {
	expr    string
	program cel.Program
	prog2   *vm.Program
}

func (ve *ValueExpr) DecodeString3(value string) (err error) {
	var ast *cel.Ast
	var env *cel.Env

	_variables := [10]*expr_proto.Decl{}
	variables := _variables[:0]

	for {
		env, err = cel.NewEnv(cel.Declarations(variables...))
		if err != nil {
			return err
		}
		var iss *cel.Issues
		ast, iss = env.Compile(value)
		if iss.Err() != nil {
			for _, e := range iss.Errors() {
				if strings.HasPrefix(e.Message, "undeclared reference to '") {
					msg := e.Message[25:]
					msg = msg[0:strings.IndexRune(msg, '\'')]
					variables = append(variables, decls.NewVar(msg, decls.Any))
				} else {
					return iss.Err()
				}
			}
		} else {
			break
		}
	}
	prg, err := env.Program(ast)
	if err != nil {
		return err
	}

	ve.expr = value
	ve.program = prg

	return nil
}

// Testing out github.com/antonmedv/expr

func (ve *ValueExpr) DecodeString(value string) error {
	prog, err := expr.Compile(value)
	if err != nil {
		return err
	}

	ve.expr = value
	ve.prog2 = prog

	return nil
}

func (ve *ValueExpr) Expr() string {
	return ve.expr
}

func (ve *ValueExpr) Eval3(variables map[string]interface{}) (interface{}, error) {
	out, _, err := ve.program.Eval(variables)
	if err != nil {
		return nil, err
	}

	return out.Value(), nil
}

func (ve *ValueExpr) Eval(variables map[string]interface{}) (interface{}, error) {
	out, err := expr.Run(ve.prog2, variables)
	if err != nil {
		return nil, err
	}

	return out, nil
}
