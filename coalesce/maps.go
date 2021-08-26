package coalesce

import "fmt"

func ToMapSI(v interface{}) (map[string]interface{}, bool) {
	switch t := v.(type) {
	case map[interface{}]interface{}:
		return MapIItoSI(t), true
	case map[string]interface{}:
		for k, v := range t {
			t[k] = ValueIItoSI(v)
		}
		return t, true
	case map[string]string:
		return MapSStoSI(t), true
	}

	return nil, false
}

func MapIItoSI(m map[interface{}]interface{}) map[string]interface{} {
	ret := make(map[string]interface{}, len(m))
	for k, v := range m {
		v = ValueIItoSI(v)
		ret[interfaceToString(k)] = v
	}
	return ret
}

func ValueIItoSI(value interface{}) interface{} {
	switch t := value.(type) {
	case map[interface{}]interface{}:
		value = MapIItoSI(t)
	case map[string]string:
		value = MapSStoSI(t)
	case []interface{}:
		for i := range t {
			t[i] = ValueIItoSI(t[i])
		}
	}
	return value
}

func MapSStoSI(m map[string]string) map[string]interface{} {
	ret := make(map[string]interface{}, len(m))
	for k, v := range m {
		ret[interfaceToString(k)] = v
	}
	return ret
}

func interfaceToString(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	if s, ok := value.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", value)
}
