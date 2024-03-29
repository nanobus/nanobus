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

type UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnum int

const (
	UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumBool UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnum = 1

	UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumMapCloudEventData UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnum = 2

	UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumArrayCloudEventData UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnum = 3

	UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumDouble UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnum = 4

	UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumString UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnum = 5
)

type UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString struct {
	Null                *types.NullVal
	Bool                bool
	MapCloudEventData   map[string]CloudEventData
	ArrayCloudEventData []CloudEventData
	Double              float64
	String              string
	UnionType           UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnum
}

func writeUnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString(r *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumBool:
		return vm.WriteBool(r.Bool, w)
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumMapCloudEventData:
		return writeMapCloudEventData(r.MapCloudEventData, w)
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumArrayCloudEventData:
		return writeArrayCloudEventData(r.ArrayCloudEventData, w)
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumDouble:
		return vm.WriteDouble(r.Double, w)
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumString:
		return vm.WriteString(r.String, w)
	}
	return fmt.Errorf("invalid value for *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString")
}

func NewUnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString() *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString {
	return &UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString{}
}

func (r *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) Serialize(w io.Writer) error {
	return writeUnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString(r, w)
}

func DeserializeUnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString(r io.Reader) (*UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString, error) {
	t := NewUnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString()
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

func DeserializeUnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringFromSchema(r io.Reader, schema string) (*UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString, error) {
	t := NewUnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString()
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

func (r *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) Schema() string {
	return "[\"null\",\"boolean\",{\"type\":\"map\",\"values\":{\"doc\":\"Representation of a JSON Value\",\"fields\":[{\"name\":\"value\",\"type\":{\"type\":\"map\",\"values\":[\"null\",\"boolean\",{\"type\":\"map\",\"values\":\"io.cloudevents.CloudEventData\"},{\"items\":\"io.cloudevents.CloudEventData\",\"type\":\"array\"},\"double\",\"string\"]}}],\"name\":\"CloudEventData\",\"type\":\"record\"}},{\"items\":\"io.cloudevents.CloudEventData\",\"type\":\"array\"},\"double\",\"string\"]"
}

func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) SetBoolean(v bool) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) SetInt(v int32) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) SetFloat(v float32) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) SetDouble(v float64) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) SetBytes(v []byte) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) SetString(v string) {
	panic("Unsupported operation")
}

func (r *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) SetLong(v int64) {

	r.UnionType = (UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnum)(v)
}

func (r *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		return &types.Boolean{Target: (&r.Bool)}
	case 2:
		r.MapCloudEventData = make(map[string]CloudEventData)
		return &MapCloudEventDataWrapper{Target: (&r.MapCloudEventData)}
	case 3:
		r.ArrayCloudEventData = make([]CloudEventData, 0)
		return &ArrayCloudEventDataWrapper{Target: (&r.ArrayCloudEventData)}
	case 4:
		return &types.Double{Target: (&r.Double)}
	case 5:
		return &types.String{Target: (&r.String)}
	}
	panic("Unknown field index")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) NullField(i int) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) HintSize(i int) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) SetDefault(i int) {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) AppendArray() types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) Finalize() {}

func (r *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumBool:
		return json.Marshal(map[string]interface{}{"boolean": r.Bool})
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumMapCloudEventData:
		return json.Marshal(map[string]interface{}{"map": r.MapCloudEventData})
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumArrayCloudEventData:
		return json.Marshal(map[string]interface{}{"array": r.ArrayCloudEventData})
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumDouble:
		return json.Marshal(map[string]interface{}{"double": r.Double})
	case UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleStringTypeEnumString:
		return json.Marshal(map[string]interface{}{"string": r.String})
	}
	return nil, fmt.Errorf("invalid value for *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString")
}

func (r *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString) UnmarshalJSON(data []byte) error {

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
	if value, ok := fields["map"]; ok {
		r.UnionType = 2
		return json.Unmarshal([]byte(value), &r.MapCloudEventData)
	}
	if value, ok := fields["array"]; ok {
		r.UnionType = 3
		return json.Unmarshal([]byte(value), &r.ArrayCloudEventData)
	}
	if value, ok := fields["double"]; ok {
		r.UnionType = 4
		return json.Unmarshal([]byte(value), &r.Double)
	}
	if value, ok := fields["string"]; ok {
		r.UnionType = 5
		return json.Unmarshal([]byte(value), &r.String)
	}
	return fmt.Errorf("invalid value for *UnionNullBoolMapCloudEventDataArrayCloudEventDataDoubleString")
}
