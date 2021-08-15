//go:generate $GOPATH/bin/gogen-avro . spec.avsc
package avro

import (
	"bytes"
	"encoding/json"

	"github.com/actgardner/gogen-avro/v9/compiler"
	"github.com/actgardner/gogen-avro/v9/vm"

	"github.com/nanobus/nanobus/codec"
	"github.com/nanobus/nanobus/resolve"
)

type (
	// Codec encodes and decodes Avro records.
	Codec struct {
		deser *vm.Program
	}
)

// JSON is the NamedLoader for this codec.
func CloudEventsAvro() (string, codec.Loader) {
	return "cloudevents+avro", Loader
}

func Loader(with interface{}, resolver resolve.ResolveAs) (codec.Codec, error) {
	t := AvroCloudEvent{}
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return NewCodec(deser), nil
}

// NewCodec creates a `Codec`.
func NewCodec(deser *vm.Program) *Codec {
	return &Codec{
		deser: deser,
	}
}

func (c *Codec) ContentType() string {
	return "application/avro"
}

// Decode decodes CloudEvents Avro bytes to a value.
func (c *Codec) Decode(msgValue []byte, args ...interface{}) (interface{}, string, error) {
	t := NewAvroCloudEvent()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, "", err
	}

	r := bytes.NewReader(msgValue)

	err = vm.Eval(r, deser, &t)
	if err != nil {
		return nil, "", err
	}

	var eventType string

	event := make(map[string]interface{}, len(t.Attribute)+1)
	for key, value := range t.Attribute {
		var v interface{}
		switch value.UnionType {
		case UnionNullBoolIntStringBytesTypeEnumBool:
			v = value.Bool
		case UnionNullBoolIntStringBytesTypeEnumInt:
			v = value.Int
		case UnionNullBoolIntStringBytesTypeEnumString:
			v = value.String
		case UnionNullBoolIntStringBytesTypeEnumBytes:
			v = value.Bytes
		}
		event[key] = v
	}

	if typeI, ok := event["type"]; ok {
		eventType, _ = typeI.(string)
	}

	switch t.Data.UnionType {
	case UnionBytesNullBoolMapUnionNullBoolAvroCloudEventDataDoubleStringArrayAvroCloudEventDataDoubleStringTypeEnumBytes:
		event["data"] = t.Data.Bytes
	case UnionBytesNullBoolMapUnionNullBoolAvroCloudEventDataDoubleStringArrayAvroCloudEventDataDoubleStringTypeEnumBool:
		event["data"] = t.Data.Bool
	case UnionBytesNullBoolMapUnionNullBoolAvroCloudEventDataDoubleStringArrayAvroCloudEventDataDoubleStringTypeEnumMapUnionNullBoolAvroCloudEventDataDoubleString:
		event["data"] = decodeDataMap(t.Data.MapUnionNullBoolAvroCloudEventDataDoubleString)
	case UnionBytesNullBoolMapUnionNullBoolAvroCloudEventDataDoubleStringArrayAvroCloudEventDataDoubleStringTypeEnumArrayAvroCloudEventData:
		event["data"] = decodeCloudEventDataArray(t.Data.ArrayAvroCloudEventData)
	case UnionBytesNullBoolMapUnionNullBoolAvroCloudEventDataDoubleStringArrayAvroCloudEventDataDoubleStringTypeEnumDouble:
		event["data"] = t.Data.Double
	case UnionBytesNullBoolMapUnionNullBoolAvroCloudEventDataDoubleStringArrayAvroCloudEventDataDoubleStringTypeEnumString:
		event["data"] = t.Data.String
	}

	return event, eventType, err
}

func decodeDataMap(d map[string]*UnionNullBoolAvroCloudEventDataDoubleString) map[string]interface{} {
	m := make(map[string]interface{}, len(d))
	for key, value := range d {
		switch value.UnionType {
		case UnionNullBoolAvroCloudEventDataDoubleStringTypeEnumBool:
			m[key] = value.Bool
		case UnionNullBoolAvroCloudEventDataDoubleStringTypeEnumAvroCloudEventData:
			m[key] = decodeCloudEventData(&value.AvroCloudEventData)
		case UnionNullBoolAvroCloudEventDataDoubleStringTypeEnumDouble:
			m[key] = value.Double
		case UnionNullBoolAvroCloudEventDataDoubleStringTypeEnumString:
			m[key] = value.String
		}
	}
	return m
}

func decodeCloudEventData(d *AvroCloudEventData) map[string]interface{} {
	m := make(map[string]interface{}, len(d.Value))
	for key, value := range d.Value {
		switch value.UnionType {
		case UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringTypeEnumBool:
			m[key] = value.Bool
		case UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringTypeEnumMapAvroCloudEventData:
			m[key] = decodeCloudEventDataMap(value.MapAvroCloudEventData)
		case UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringTypeEnumArrayAvroCloudEventData:
			m[key] = decodeCloudEventDataArray(value.ArrayAvroCloudEventData)
		case UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringTypeEnumDouble:
			m[key] = value.Double
		case UnionNullBoolMapAvroCloudEventDataArrayAvroCloudEventDataDoubleStringTypeEnumString:
			m[key] = value.String
		}
	}
	return m
}

func decodeCloudEventDataArray(d []AvroCloudEventData) []interface{} {
	m := make([]interface{}, len(d))
	for i := range d {
		value := &d[i]
		m[i] = decodeCloudEventData(value)
	}
	return m
}

func decodeCloudEventDataMap(d map[string]AvroCloudEventData) map[string]interface{} {
	m := make(map[string]interface{}, len(d))
	for key, value := range d {
		m[key] = decodeCloudEventData(&value)
	}
	return m
}

// Encode encodes a value into CloudEvents Avro encoded bytes.
func (c *Codec) Encode(value interface{}, args ...interface{}) ([]byte, error) {
	return json.Marshal(value)
}
