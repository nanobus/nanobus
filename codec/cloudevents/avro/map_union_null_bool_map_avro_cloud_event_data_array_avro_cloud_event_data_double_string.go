// Code generated by github.com/actgardner/gogen-avro/v8. DO NOT EDIT.
/*
 * SOURCE:
 *     spec.avsc
 */
package avro

import (
	"github.com/actgardner/gogen-avro/v9/vm"
	"github.com/actgardner/gogen-avro/v9/vm/types"
	"io"
)

func writeMapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleString(r map[string]*UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleString, w io.Writer) error {
	err := vm.WriteLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
		return err
	}
	for k, e := range r {
		err = vm.WriteString(k, w)
		if err != nil {
			return err
		}
		err = writeUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleString(e, w)
		if err != nil {
			return err
		}
	}
	return vm.WriteLong(0, w)
}

type MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper struct {
	Target *map[string]*UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleString
	keys   []string
	values []*UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleString
}

func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetBoolean(v bool) {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetInt(v int32) {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetLong(v int64) {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetFloat(v float32) {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetDouble(v float64) {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetBytes(v []byte) {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetString(v string) {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetUnionElem(v int64) {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) Get(i int) types.Field {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) SetDefault(i int) {
	panic("Unsupported operation")
}

func (r *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) NullField(_ int) {
	r.values[len(r.values)-1] = nil
}

func (r *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) Finalize() {
	for i := range r.keys {
		(*r.Target)[r.keys[i]] = r.values[i]
	}
}

func (r *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) AppendMap(key string) types.Field {
	r.keys = append(r.keys, key)
	var v *UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleString
	v = NewUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleString()

	r.values = append(r.values, v)
	return r.values[len(r.values)-1]
}

func (_ *MapUnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringWrapper) AppendArray() types.Field {
	panic("Unsupported operation")
}
