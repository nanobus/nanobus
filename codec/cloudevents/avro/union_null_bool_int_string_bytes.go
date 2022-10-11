// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCE:
 *     spec.avsc
 */
package avro

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
)

type UnionNullBoolIntStringBytesTypeEnum int

const (
	UnionNullBoolIntStringBytesTypeEnumBool UnionNullBoolIntStringBytesTypeEnum = 1

	UnionNullBoolIntStringBytesTypeEnumInt UnionNullBoolIntStringBytesTypeEnum = 2

	UnionNullBoolIntStringBytesTypeEnumString UnionNullBoolIntStringBytesTypeEnum = 3

	UnionNullBoolIntStringBytesTypeEnumBytes UnionNullBoolIntStringBytesTypeEnum = 4
)

type UnionNullBoolIntStringBytes struct {
	Null      *types.NullVal
	Bool      bool
	Int       int32
	String    string
	Bytes     Bytes
	UnionType UnionNullBoolIntStringBytesTypeEnum
}

func writeUnionNullBoolIntStringBytes(r *UnionNullBoolIntStringBytes, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullBoolIntStringBytesTypeEnumBool:
		return vm.WriteBool(r.Bool, w)
	case UnionNullBoolIntStringBytesTypeEnumInt:
		return vm.WriteInt(r.Int, w)
	case UnionNullBoolIntStringBytesTypeEnumString:
		return vm.WriteString(r.String, w)
	case UnionNullBoolIntStringBytesTypeEnumBytes:
		return vm.WriteBytes(r.Bytes, w)
	}
	return fmt.Errorf("invalid value for *UnionNullBoolIntStringBytes")
}

func NewUnionNullBoolIntStringBytes() *UnionNullBoolIntStringBytes {
	return &UnionNullBoolIntStringBytes{}
}

func (r *UnionNullBoolIntStringBytes) Serialize(w io.Writer) error {
	return writeUnionNullBoolIntStringBytes(r, w)
}

func DeserializeUnionNullBoolIntStringBytes(r io.Reader) (*UnionNullBoolIntStringBytes, error) {
	t := NewUnionNullBoolIntStringBytes()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func DeserializeUnionNullBoolIntStringBytesFromSchema(r io.Reader, schema string) (*UnionNullBoolIntStringBytes, error) {
	t := NewUnionNullBoolIntStringBytes()
	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func (r *UnionNullBoolIntStringBytes) Schema() string {
	return "[\"null\",\"boolean\",\"int\",\"string\",\"bytes\"]"
}

func (_ *UnionNullBoolIntStringBytes) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) SetString(v string)  { panic("Unsupported operation") }

func (r *UnionNullBoolIntStringBytes) SetLong(v int64) {

	r.UnionType = (UnionNullBoolIntStringBytesTypeEnum)(v)
}

func (r *UnionNullBoolIntStringBytes) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		return &types.Boolean{Target: (&r.Bool)}
	case 2:
		return &types.Int{Target: (&r.Int)}
	case 3:
		return &types.String{Target: (&r.String)}
	case 4:
		return &BytesWrapper{Target: (&r.Bytes)}
	}
	panic("Unknown field index")
}
func (_ *UnionNullBoolIntStringBytes) NullField(i int)  { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) HintSize(i int)   { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) SetDefault(i int) { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolIntStringBytes) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *UnionNullBoolIntStringBytes) Finalize()                {}

func (r *UnionNullBoolIntStringBytes) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionNullBoolIntStringBytesTypeEnumBool:
		return json.Marshal(map[string]interface{}{"boolean": r.Bool})
	case UnionNullBoolIntStringBytesTypeEnumInt:
		return json.Marshal(map[string]interface{}{"int": r.Int})
	case UnionNullBoolIntStringBytesTypeEnumString:
		return json.Marshal(map[string]interface{}{"string": r.String})
	case UnionNullBoolIntStringBytesTypeEnumBytes:
		return json.Marshal(map[string]interface{}{"bytes": r.Bytes})
	}
	return nil, fmt.Errorf("invalid value for *UnionNullBoolIntStringBytes")
}

func (r *UnionNullBoolIntStringBytes) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["boolean"]; ok {
		r.UnionType = 1
		return json.Unmarshal([]byte(value), &r.Bool)
	}
	if value, ok := fields["int"]; ok {
		r.UnionType = 2
		return json.Unmarshal([]byte(value), &r.Int)
	}
	if value, ok := fields["string"]; ok {
		r.UnionType = 3
		return json.Unmarshal([]byte(value), &r.String)
	}
	if value, ok := fields["bytes"]; ok {
		r.UnionType = 4
		return json.Unmarshal([]byte(value), &r.Bytes)
	}
	return fmt.Errorf("invalid value for *UnionNullBoolIntStringBytes")
}
