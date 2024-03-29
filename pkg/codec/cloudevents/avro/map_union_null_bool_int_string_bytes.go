// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCE:
 *     spec.avsc
 */
package avro

import (
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
	"io"
)

func writeMapUnionNullBoolIntStringBytes(r map[string]*UnionNullBoolIntStringBytes, w io.Writer) error {
	err := vm.WriteLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
		return err
	}
	for k, e := range r {
		err = vm.WriteString(k, w)
		if err != nil {
			return err
		}
		err = writeUnionNullBoolIntStringBytes(e, w)
		if err != nil {
			return err
		}
	}
	return vm.WriteLong(0, w)
}

type MapUnionNullBoolIntStringBytesWrapper struct {
	Target *map[string]*UnionNullBoolIntStringBytes
	keys   []string
	values []*UnionNullBoolIntStringBytes
}

func (_ *MapUnionNullBoolIntStringBytesWrapper) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *MapUnionNullBoolIntStringBytesWrapper) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *MapUnionNullBoolIntStringBytesWrapper) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *MapUnionNullBoolIntStringBytesWrapper) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *MapUnionNullBoolIntStringBytesWrapper) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *MapUnionNullBoolIntStringBytesWrapper) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *MapUnionNullBoolIntStringBytesWrapper) SetString(v string)   { panic("Unsupported operation") }
func (_ *MapUnionNullBoolIntStringBytesWrapper) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *MapUnionNullBoolIntStringBytesWrapper) Get(i int) types.Field {
	panic("Unsupported operation")
}
func (_ *MapUnionNullBoolIntStringBytesWrapper) SetDefault(i int) { panic("Unsupported operation") }

func (r *MapUnionNullBoolIntStringBytesWrapper) HintSize(s int) {
	if r.keys == nil {
		r.keys = make([]string, 0, s)
		r.values = make([]*UnionNullBoolIntStringBytes, 0, s)
	}
}

func (r *MapUnionNullBoolIntStringBytesWrapper) NullField(_ int) {
	r.values[len(r.values)-1] = nil
}

func (r *MapUnionNullBoolIntStringBytesWrapper) Finalize() {
	for i := range r.keys {
		(*r.Target)[r.keys[i]] = r.values[i]
	}
}

func (r *MapUnionNullBoolIntStringBytesWrapper) AppendMap(key string) types.Field {
	r.keys = append(r.keys, key)
	var v *UnionNullBoolIntStringBytes
	v = NewUnionNullBoolIntStringBytes()

	r.values = append(r.values, v)
	return r.values[len(r.values)-1]
}

func (_ *MapUnionNullBoolIntStringBytesWrapper) AppendArray() types.Field {
	panic("Unsupported operation")
}
