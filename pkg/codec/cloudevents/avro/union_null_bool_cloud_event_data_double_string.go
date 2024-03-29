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

type UnionNullBoolCloudEventDataDoubleStringTypeEnum int

const (
	UnionNullBoolCloudEventDataDoubleStringTypeEnumBool UnionNullBoolCloudEventDataDoubleStringTypeEnum = 1

	UnionNullBoolCloudEventDataDoubleStringTypeEnumCloudEventData UnionNullBoolCloudEventDataDoubleStringTypeEnum = 2

	UnionNullBoolCloudEventDataDoubleStringTypeEnumDouble UnionNullBoolCloudEventDataDoubleStringTypeEnum = 3

	UnionNullBoolCloudEventDataDoubleStringTypeEnumString UnionNullBoolCloudEventDataDoubleStringTypeEnum = 4
)

type UnionNullBoolCloudEventDataDoubleString struct {
	Null           *types.NullVal
	Bool           bool
	CloudEventData CloudEventData
	Double         float64
	String         string
	UnionType      UnionNullBoolCloudEventDataDoubleStringTypeEnum
}

func writeUnionNullBoolCloudEventDataDoubleString(r *UnionNullBoolCloudEventDataDoubleString, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullBoolCloudEventDataDoubleStringTypeEnumBool:
		return vm.WriteBool(r.Bool, w)
	case UnionNullBoolCloudEventDataDoubleStringTypeEnumCloudEventData:
		return writeCloudEventData(r.CloudEventData, w)
	case UnionNullBoolCloudEventDataDoubleStringTypeEnumDouble:
		return vm.WriteDouble(r.Double, w)
	case UnionNullBoolCloudEventDataDoubleStringTypeEnumString:
		return vm.WriteString(r.String, w)
	}
	return fmt.Errorf("invalid value for *UnionNullBoolCloudEventDataDoubleString")
}

func NewUnionNullBoolCloudEventDataDoubleString() *UnionNullBoolCloudEventDataDoubleString {
	return &UnionNullBoolCloudEventDataDoubleString{}
}

func (r *UnionNullBoolCloudEventDataDoubleString) Serialize(w io.Writer) error {
	return writeUnionNullBoolCloudEventDataDoubleString(r, w)
}

func DeserializeUnionNullBoolCloudEventDataDoubleString(r io.Reader) (*UnionNullBoolCloudEventDataDoubleString, error) {
	t := NewUnionNullBoolCloudEventDataDoubleString()
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

func DeserializeUnionNullBoolCloudEventDataDoubleStringFromSchema(r io.Reader, schema string) (*UnionNullBoolCloudEventDataDoubleString, error) {
	t := NewUnionNullBoolCloudEventDataDoubleString()
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

func (r *UnionNullBoolCloudEventDataDoubleString) Schema() string {
	return "[\"null\",\"boolean\",{\"doc\":\"Representation of a JSON Value\",\"fields\":[{\"name\":\"value\",\"type\":{\"type\":\"map\",\"values\":[\"null\",\"boolean\",{\"type\":\"map\",\"values\":\"io.cloudevents.CloudEventData\"},{\"items\":\"io.cloudevents.CloudEventData\",\"type\":\"array\"},\"double\",\"string\"]}}],\"name\":\"CloudEventData\",\"type\":\"record\"},\"double\",\"string\"]"
}

func (_ *UnionNullBoolCloudEventDataDoubleString) SetBoolean(v bool)  { panic("Unsupported operation") }
func (_ *UnionNullBoolCloudEventDataDoubleString) SetInt(v int32)     { panic("Unsupported operation") }
func (_ *UnionNullBoolCloudEventDataDoubleString) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *UnionNullBoolCloudEventDataDoubleString) SetDouble(v float64) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolCloudEventDataDoubleString) SetBytes(v []byte)  { panic("Unsupported operation") }
func (_ *UnionNullBoolCloudEventDataDoubleString) SetString(v string) { panic("Unsupported operation") }

func (r *UnionNullBoolCloudEventDataDoubleString) SetLong(v int64) {

	r.UnionType = (UnionNullBoolCloudEventDataDoubleStringTypeEnum)(v)
}

func (r *UnionNullBoolCloudEventDataDoubleString) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		return &types.Boolean{Target: (&r.Bool)}
	case 2:
		r.CloudEventData = NewCloudEventData()
		return &types.Record{Target: (&r.CloudEventData)}
	case 3:
		return &types.Double{Target: (&r.Double)}
	case 4:
		return &types.String{Target: (&r.String)}
	}
	panic("Unknown field index")
}
func (_ *UnionNullBoolCloudEventDataDoubleString) NullField(i int)  { panic("Unsupported operation") }
func (_ *UnionNullBoolCloudEventDataDoubleString) HintSize(i int)   { panic("Unsupported operation") }
func (_ *UnionNullBoolCloudEventDataDoubleString) SetDefault(i int) { panic("Unsupported operation") }
func (_ *UnionNullBoolCloudEventDataDoubleString) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolCloudEventDataDoubleString) AppendArray() types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolCloudEventDataDoubleString) Finalize() {}

func (r *UnionNullBoolCloudEventDataDoubleString) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionNullBoolCloudEventDataDoubleStringTypeEnumBool:
		return json.Marshal(map[string]interface{}{"boolean": r.Bool})
	case UnionNullBoolCloudEventDataDoubleStringTypeEnumCloudEventData:
		return json.Marshal(map[string]interface{}{"io.cloudevents.CloudEventData": r.CloudEventData})
	case UnionNullBoolCloudEventDataDoubleStringTypeEnumDouble:
		return json.Marshal(map[string]interface{}{"double": r.Double})
	case UnionNullBoolCloudEventDataDoubleStringTypeEnumString:
		return json.Marshal(map[string]interface{}{"string": r.String})
	}
	return nil, fmt.Errorf("invalid value for *UnionNullBoolCloudEventDataDoubleString")
}

func (r *UnionNullBoolCloudEventDataDoubleString) UnmarshalJSON(data []byte) error {

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
	if value, ok := fields["io.cloudevents.CloudEventData"]; ok {
		r.UnionType = 2
		return json.Unmarshal([]byte(value), &r.CloudEventData)
	}
	if value, ok := fields["double"]; ok {
		r.UnionType = 3
		return json.Unmarshal([]byte(value), &r.Double)
	}
	if value, ok := fields["string"]; ok {
		r.UnionType = 4
		return json.Unmarshal([]byte(value), &r.String)
	}
	return fmt.Errorf("invalid value for *UnionNullBoolCloudEventDataDoubleString")
}
