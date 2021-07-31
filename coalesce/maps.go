package coalesce

import "fmt"

func ToMapSI(v interface{}) (map[string]interface{}, bool) {
	switch t := v.(type) {
	case map[interface{}]interface{}:
		return MapIItoSI(t), true
	case map[string]interface{}:
		return t, true
	case map[string]string:
		return MapSStoSI(t), true
	}

	return nil, false
}

func MapIItoSI(m map[interface{}]interface{}) map[string]interface{} {
	ret := make(map[string]interface{}, len(m))
	for k, v := range m {
		ret[interfaceToString(k)] = v
	}
	return ret
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
