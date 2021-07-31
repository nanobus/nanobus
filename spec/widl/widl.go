package widl

import (
	"os"

	"github.com/wapc/widl-go/ast"
	"github.com/wapc/widl-go/parser"

	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/spec"
)

type WIDLConfig struct {
	// FileName is the file name of the WIDL definition to load.
	FileName string `mapstructure:"fileName"` // TODO: Load from external location
}

// WIDL is the NamedLoader for the WIDL spec.
func WIDL() (string, spec.Loader) {
	return "widl", WIDLLoader
}

func WIDLLoader(with interface{}) ([]*spec.Namespace, error) {
	c := WIDLConfig{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	widlBytes, err := os.ReadFile(c.FileName)
	if err != nil {
		return nil, err
	}

	ns, err := Parse(widlBytes)
	if err != nil {
		return nil, err
	}

	return []*spec.Namespace{ns}, nil
}

type nsParser struct {
	n *spec.Namespace
}

func Parse(schema []byte) (*spec.Namespace, error) {
	n := spec.Namespace{
		Services: make(map[string]*spec.Service),
		Types:    make(map[string]*spec.Type),
		Enums:    make(map[string]*spec.Enum),
		Unions:   make(map[string]*spec.Union),
	}
	doc, err := parser.Parse(parser.ParseParams{
		Source: string(schema),
		Options: parser.ParseOptions{
			NoLocation: true,
			NoSource:   true,
		},
	})
	if err != nil {
		return nil, err
	}

	p := nsParser{n: &n}

	for _, def := range doc.Definitions {
		switch d := def.(type) {
		case *ast.NamespaceDefinition:
			n.Name = d.Name.Value
			n.Annotated = p.convertAnnotations(d.Annotations)

		case *ast.TypeDefinition:
			t := p.createType(d)
			n.Types[t.Name] = t
		}
	}

	for _, def := range doc.Definitions {
		switch d := def.(type) {
		case *ast.RoleDefinition:
			//if a := d.Annotation("service"); a != nil {
			s := p.convertService(d)
			n.Services[s.Name] = s
			//}

		case *ast.TypeDefinition:
			p.convertType(d)
		}
	}

	return &n, nil
}

func (p *nsParser) createType(tt *ast.TypeDefinition) *spec.Type {
	t := spec.Type{
		Name: tt.Name.Value,
	}
	return &t
}

func (p *nsParser) convertType(tt *ast.TypeDefinition) {
	t := p.n.Types[tt.Name.Value]
	*t = spec.Type{
		Name:      tt.Name.Value,
		Fields:    p.convertFields(tt.Fields),
		Annotated: p.convertAnnotations(tt.Annotations),
	}
}

func (p *nsParser) convertFields(fields []*ast.FieldDefinition) map[string]*spec.Field {
	if fields == nil {
		return nil
	}

	o := make(map[string]*spec.Field, len(fields))
	for _, field := range fields {
		var dv interface{}
		if field.Default != nil {
			dv = field.Default.GetValue()
		}
		o[field.Name.Value] = &spec.Field{
			Name:         field.Name.Value,
			Type:         p.convertTypeRef(field.Type),
			DefaultValue: dv,
			Annotated:    p.convertAnnotations(field.Annotations),
		}
	}

	return o
}

func (p *nsParser) convertService(role *ast.RoleDefinition) *spec.Service {
	s := spec.Service{
		Name:       role.Name.Value,
		Operations: p.convertOperations(role.Operations),
		Annotated:  p.convertAnnotations(role.Annotations),
	}
	return &s
}

func (p *nsParser) convertOperations(operations []*ast.OperationDefinition) map[string]*spec.Operation {
	if operations == nil {
		return nil
	}

	o := make(map[string]*spec.Operation, len(operations))
	for _, operation := range operations {
		var params *spec.Type
		if operation.Unary {
			if named, ok := operation.Parameters[0].Type.(*ast.Named); ok {
				params = p.n.Types[named.Name.Value]
			}
		} else {
			params = p.convertParameterType(operation.Name.Value+"Params", operation.Parameters)
		}
		o[operation.Name.Value] = &spec.Operation{
			Name:       operation.Name.Value,
			Unary:      operation.Unary,
			Parameters: params,
			Returns:    p.convertTypeRef(operation.Type),
			Annotated:  p.convertAnnotations(operation.Annotations),
		}
	}

	return o
}

func (p *nsParser) convertParameterType(name string, params []*ast.ParameterDefinition) *spec.Type {
	t := spec.Type{
		Namespace: p.n,
		Name:      name,
		Fields:    p.convertParameters(params),
	}
	return &t
}

func (p *nsParser) convertParameters(parameters []*ast.ParameterDefinition) map[string]*spec.Field {
	if parameters == nil {
		return nil
	}

	o := make(map[string]*spec.Field, len(parameters))
	for _, parameter := range parameters {
		var dv interface{}
		if parameter.Default != nil {
			dv = parameter.Default.GetValue()
		}
		o[parameter.Name.Value] = &spec.Field{
			Name:         parameter.Name.Value,
			Type:         p.convertTypeRef(parameter.Type),
			DefaultValue: dv,
			Annotated:    p.convertAnnotations(parameter.Annotations),
		}
	}

	return o
}

func (p *nsParser) convertAnnotations(annotations []*ast.Annotation) spec.Annotated {
	if annotations == nil {
		return spec.Annotated{}
	}

	a := make(map[string]*spec.Annotation, len(annotations))
	for _, annotation := range annotations {
		a[annotation.Name.Value] = &spec.Annotation{
			Name:      annotation.Name.Value,
			Arguments: p.convertArguments(annotation.Arguments),
		}
	}

	return spec.Annotated{
		Annotations: a,
	}
}

func (p *nsParser) convertArguments(arguments []*ast.Argument) map[string]*spec.Argument {
	if arguments == nil {
		return nil
	}

	a := make(map[string]*spec.Argument, len(arguments))
	for _, argument := range arguments {
		a[argument.Name.Value] = &spec.Argument{
			Name:  argument.Name.Value,
			Value: argument.Value.GetValue(),
		}
	}

	return a
}

var (
	typeRefString   = spec.TypeRef{Kind: spec.KindString}
	typeRefU64      = spec.TypeRef{Kind: spec.KindU64}
	typeRefU32      = spec.TypeRef{Kind: spec.KindU32}
	typeRefU16      = spec.TypeRef{Kind: spec.KindU16}
	typeRefU8       = spec.TypeRef{Kind: spec.KindU8}
	typeRefI64      = spec.TypeRef{Kind: spec.KindI64}
	typeRefI32      = spec.TypeRef{Kind: spec.KindI32}
	typeRefI16      = spec.TypeRef{Kind: spec.KindI16}
	typeRefI8       = spec.TypeRef{Kind: spec.KindI8}
	typeRefF64      = spec.TypeRef{Kind: spec.KindF64}
	typeRefF32      = spec.TypeRef{Kind: spec.KindF32}
	typeRefBool     = spec.TypeRef{Kind: spec.KindBool}
	typeRefBytes    = spec.TypeRef{Kind: spec.KindBytes}
	typeRefRaw      = spec.TypeRef{Kind: spec.KindRaw}
	typeRefDateTime = spec.TypeRef{Kind: spec.KindDateTime}

	typeRefMap = map[string]*spec.TypeRef{
		"string":   &typeRefString,
		"u64":      &typeRefU64,
		"u32":      &typeRefU32,
		"u16":      &typeRefU16,
		"u8":       &typeRefU8,
		"i64":      &typeRefI64,
		"i32":      &typeRefI32,
		"i16":      &typeRefI16,
		"i8":       &typeRefI8,
		"f64":      &typeRefF64,
		"f32":      &typeRefF32,
		"bool":     &typeRefBool,
		"bytes":    &typeRefBytes,
		"raw":      &typeRefRaw,
		"datetime": &typeRefDateTime,
	}
)

func (p *nsParser) convertTypeRef(t ast.Type) *spec.TypeRef {
	if t == nil {
		return nil
	}

	switch tt := t.(type) {
	case *ast.Named:
		if tt.Name.Value == "void" {
			return nil
		}
		if prim, ok := typeRefMap[tt.Name.Value]; ok {
			return prim
		}
		if t, ok := p.n.Types[tt.Name.Value]; ok {
			return &spec.TypeRef{
				Kind: spec.KindType,
				Type: t,
			}
		}
		if e, ok := p.n.Enums[tt.Name.Value]; ok {
			return &spec.TypeRef{
				Kind: spec.KindEnum,
				Enum: e,
			}
		}
		if u, ok := p.n.Unions[tt.Name.Value]; ok {
			return &spec.TypeRef{
				Kind:  spec.KindUnion,
				Union: u,
			}
		}
	case *ast.ListType:
		return &spec.TypeRef{
			Kind:     spec.KindList,
			ListType: p.convertTypeRef(tt.Type),
		}
	case *ast.MapType:
		return &spec.TypeRef{
			Kind:         spec.KindMap,
			MapKeyType:   p.convertTypeRef(tt.KeyType),
			MapValueType: p.convertTypeRef(tt.ValueType),
		}
	case *ast.Optional:
		return &spec.TypeRef{
			Kind:         spec.KindOptional,
			OptionalType: p.convertTypeRef(tt.Type),
		}
	}

	panic("unreachable")
}
