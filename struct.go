package goutils

import (
	"reflect"
	"strings"
)

// Detect the StructField is a number field, isn't it.
func IsNumberField(f reflect.StructField) bool {
	switch f.Type.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		return true
	default:
		return false
	}
}

// Get json field name from StructTag
func GetJSONTag(t reflect.StructTag) string {
	jsonTag := t.Get("json")
	if jsonTag == "" {
		return ""
	}
	return strings.Split(jsonTag, ",")[0]
}
