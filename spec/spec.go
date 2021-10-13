package spec

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"

	"github.com/nanobus/nanobus/coalesce"
)

type (
	NamedLoader func() (string, Loader)
	Loader      func(config interface{}) ([]*Namespace, error)
	Registry    map[string]Loader
)

func (r Registry) Register(loaders ...NamedLoader) {
	for _, l := range loaders {
		name, loader := l()
		r[name] = loader
	}
}

type (
	Namespaces map[string]*Namespace

	Namespace struct {
		Name           string              `json:"name"`
		Services       []*Service          `json:"services"`
		ServicesByName map[string]*Service `json:"-"`
		Types          []*Type             `json:"types"`
		TypesByName    map[string]*Type    `json:"-"`
		Enums          map[string]*Enum    `json:"enums"`
		Unions         map[string]*Union   `json:"unions"`
		Annotated
	}

	Service struct {
		Name             string                `json:"name"`
		Description      string                `json:"description,omitempty"`
		Operations       []*Operation          `json:"operations"`
		OperationsByName map[string]*Operation `json:"-"`
		Annotated
	}

	Operation struct {
		Name        string   `json:"name"`
		Description string   `json:"description,omitempty"`
		Unary       bool     `json:"unary"`
		Parameters  *Type    `json:"parameters,omitempty"`
		Returns     *TypeRef `json:"returns,omitempty"`
		Annotated
	}

	Type struct {
		Namespace    *Namespace        `json:"-"`
		Name         string            `json:"name"`
		Description  string            `json:"description,omitempty"`
		Fields       []*Field          `json:"fields"`
		FieldsByName map[string]*Field `json:"-"`
		Annotated
		Validations []Validation `json:"-"`
	}

	Field struct {
		Name         string      `json:"name"`
		Description  string      `json:"description,omitempty"`
		Type         *TypeRef    `json:"type"`
		DefaultValue interface{} `json:"defaultValue,omitempty"`
		Annotated
	}

	Enum struct {
		Namespace   *Namespace   `json:"-"`
		Name        string       `json:"name"`
		Description string       `json:"description,omitempty"`
		Values      []*EnumValue `json:"values"`
		Annotated
	}

	EnumValue struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		StringValue string `json:"stringValue"`
		IndexValue  int    `json:"indexValue"`
		Annotated
	}

	Union struct {
		Namespace   *Namespace `json:"-"`
		Name        string     `json:"name"`
		Description string     `json:"description,omitempty"`
		Types       []*TypeRef `json:"types"`
		Annotated
	}

	Annotated struct {
		Annotations map[string]*Annotation `json:"annotations,omitempty"`
	}

	Annotation struct {
		Name      string               `json:"name"`
		Arguments map[string]*Argument `json:"arguments,omitempty"`
	}

	Argument struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}

	TypeRef struct {
		Kind         Kind     `json:"kind"`
		Type         *Type    `json:"type,omitempty"`
		Enum         *Enum    `json:"enum,omitempty"`
		Union        *Union   `json:"union,omitempty"`
		OptionalType *TypeRef `json:"optionalType,omitempty"`
		ListType     *TypeRef `json:"listType,omitempty"`
		MapKeyType   *TypeRef `json:"mapKeyType,omitempty"`
		MapValueType *TypeRef `json:"mapValuetype,omitempty"`
	}

	Validation func(v interface{}) ([]ValidationError, error)

	ValidationError struct {
		Fields   []string `json:"fields"`
		Messages []string `json:"messages"`
	}

	Annotator interface {
		Annotation(name string) (*Annotation, bool)
	}
)

func (ns Namespaces) Operation(namespace, service, operation string) (*Operation, bool) {
	n, ok := ns[namespace]
	if !ok {
		return nil, false
	}
	s, ok := n.ServicesByName[service]
	if !ok {
		return nil, false
	}
	o, ok := s.OperationsByName[operation]
	return o, ok
}

func (t *TypeRef) IsPrimitive() bool {
	switch t.Kind {
	case KindOptional:
		return t.OptionalType.IsPrimitive()
	}
	return t.Kind.IsPrimitive()
}

func (t *TypeRef) Coalesce(value interface{}, validate bool) (interface{}, error) {
	var err error
	switch t.Kind {
	case KindOptional:
		if value == nil {
			return nil, nil
		}
		return t.OptionalType.Coalesce(value, validate)
	case KindString, KindDateTime:
		if _, ok := value.(string); !ok {
			err = fmt.Errorf("value must be an string")
		}
	case KindU64:
		if _, ok := value.(uint64); !ok {
			value, err = cast.ToUint64E(value)
		}
	case KindU32:
		if _, ok := value.(uint32); !ok {
			value, err = cast.ToUint32E(value)
		}
	case KindU16:
		if _, ok := value.(uint16); !ok {
			value, err = cast.ToUint16E(value)
		}
	case KindU8:
		if _, ok := value.(uint8); !ok {
			value, err = cast.ToUint8E(value)
		}
	case KindI64:
		if _, ok := value.(int64); !ok {
			value, err = cast.ToInt64E(value)
		}
	case KindI32:
		if _, ok := value.(int32); !ok {
			value, err = cast.ToInt32E(value)
		}
	case KindI16:
		if _, ok := value.(int16); !ok {
			value, err = cast.ToInt16E(value)
		}
	case KindI8:
		if _, ok := value.(int8); !ok {
			value, err = cast.ToInt8E(value)
		}
	case KindF64:
		if _, ok := value.(float64); !ok {
			value, err = cast.ToFloat64E(value)
		}
	case KindF32:
		if _, ok := value.(float32); !ok {
			value, err = cast.ToFloat32E(value)
		}
	case KindBool:
		if _, ok := value.(bool); !ok {
			value, err = cast.ToBoolE(value)
		}
	case KindBytes:
		if _, ok := value.([]byte); !ok {
			if stringValue, ok := value.(string); ok {
				value, err = base64.StdEncoding.DecodeString(stringValue)
			} else {
				err = fmt.Errorf("value must be a boolean")
			}
		}
	case KindType:
		valueMap, ok := coalesce.ToMapSI(value)
		if !ok {
			err = fmt.Errorf("value must be a map")
		}
		if err == nil {
			err = t.Type.Coalesce(valueMap, validate)
		}
		value = valueMap
		//KindEnum
		//KindUnion
	}

	return value, err
}

func (t *Type) Coalesce(v map[string]interface{}, validate bool) error {
	for fieldName, value := range v {
		f, ok := t.FieldsByName[fieldName]
		if !ok {
			// Exclude extraneous values.
			delete(v, fieldName)
			continue
		}

		if err := t.doField(f.Type, f, fieldName, v, value, validate); err != nil {
			return err
		}
	}

	if validate {
		for fieldName, f := range t.FieldsByName {
			if f.Type.Kind != KindOptional {
				if _, ok := v[fieldName]; !ok {
					return fmt.Errorf("missing required field %s in type %s", fieldName, t.Name)
				}
			}
		}
	}

	return nil
}

func (t *Type) doField(tt *TypeRef, f *Field, fieldName string, v map[string]interface{}, value interface{}, validate bool) (err error) {
	switch tt.Kind {
	case KindOptional:
		if value == nil {
			return nil
		}
		return t.doField(tt.OptionalType, f, fieldName, v, value, validate)
	case KindString, KindDateTime:
		if _, ok := value.(string); !ok {
			err = fmt.Errorf("field %q of type %q must be a string", f.Name, t.Name)
		}
	case KindU64:
		if _, ok := value.(uint64); !ok {
			v[fieldName], err = cast.ToUint64E(value)
		}
	case KindU32:
		if _, ok := value.(uint32); !ok {
			v[fieldName], err = cast.ToUint32E(value)
		}
	case KindU16:
		if _, ok := value.(uint16); !ok {
			v[fieldName], err = cast.ToUint16E(value)
		}
	case KindU8:
		if _, ok := value.(uint8); !ok {
			v[fieldName], err = cast.ToUint8E(value)
		}
	case KindI64:
		if _, ok := value.(int64); !ok {
			v[fieldName], err = cast.ToInt64E(value)
		}
	case KindI32:
		if _, ok := value.(int32); !ok {
			v[fieldName], err = cast.ToInt32E(value)
		}
	case KindI16:
		if _, ok := value.(int16); !ok {
			v[fieldName], err = cast.ToInt16E(value)
		}
	case KindI8:
		if _, ok := value.(int8); !ok {
			v[fieldName], err = cast.ToInt8E(value)
		}
	case KindF64:
		if _, ok := value.(float64); !ok {
			v[fieldName], err = cast.ToFloat64E(value)
		}
	case KindF32:
		if _, ok := value.(float32); !ok {
			v[fieldName], err = cast.ToFloat32E(value)
		}
	case KindBool:
		if _, ok := value.(bool); !ok {
			v[fieldName], err = cast.ToBoolE(value)
		}
	case KindBytes:
		if _, ok := value.([]byte); !ok {
			if stringValue, ok := value.(string); ok {
				value, err = base64.StdEncoding.DecodeString(stringValue)
			} else {
				err = fmt.Errorf("field %q of type %q must be a boolean", f.Name, t.Name)
			}
		}
	case KindType:
		valueMap, ok := coalesce.ToMapSI(value)
		if !ok {
			err = fmt.Errorf("field %q of type %q must be a map", f.Name, t.Name)
		}
		if err == nil {
			err = tt.Type.Coalesce(valueMap, validate)
		}
		//KindEnum
		//KindUnion
	}

	return err
}

func (a *Annotated) Annotation(name string) (*Annotation, bool) {
	anno, ok := a.Annotations[name]
	return anno, ok
}

type Kind int

const (
	KindUnknown Kind = iota
	KindOptional
	KindList
	KindMap
	KindString
	KindU64
	KindU32
	KindU16
	KindU8
	KindI64
	KindI32
	KindI16
	KindI8
	KindF64
	KindF32
	KindBool
	KindBytes
	KindRaw
	KindDateTime
	KindType
	KindEnum
	KindUnion
)

func (k Kind) String() string {
	switch k {
	case KindOptional:
		return "optional"
	case KindList:
		return "list"
	case KindMap:
		return "map"
	case KindString:
		return "string"
	case KindU64:
		return "u64"
	case KindU32:
		return "u32"
	case KindU16:
		return "u16"
	case KindU8:
		return "u8"
	case KindI64:
		return "i64"
	case KindI32:
		return "i32"
	case KindI16:
		return "i16"
	case KindI8:
		return "i8"
	case KindF64:
		return "f64"
	case KindF32:
		return "f32"
	case KindBool:
		return "bool"
	case KindBytes:
		return "bytes"
	case KindRaw:
		return "raw"
	case KindDateTime:
		return "datetime"
	case KindType:
		return "type"
	case KindEnum:
		return "enum"
	case KindUnion:
		return "union"
	}
	return "unknown"
}

func (k Kind) IsPrimitive() bool {
	switch k {
	case KindString, KindU64, KindU32, KindU16, KindU8, KindI64,
		KindI32, KindI16, KindI8, KindF64, KindF32, KindBool, KindDateTime:
		return true
	}
	return false
}

func (k Kind) MarshalJSON() ([]byte, error) {
	s := k.String()
	return json.Marshal(s)
}

func (a *Argument) ValueString() string {
	return cast.ToString(a.Value)
}
