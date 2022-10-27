package goutils

import (
	"github.com/goccy/go-json"
)

type JSON json.RawMessage

// Convert any object to JSON string
func Marshal(origin interface{}) string {
	bytes, err := json.Marshal(origin)
	if err != nil {
		Error(err)
	}
	return string(bytes)
}

// Another name for UnsafeConvert
func Unmarshal[T any](origin string) T {
	return UnsafeConvert[T](origin)
}

// Convert a string or struct to another struct.
// This func uses JSON as a middle data type to convert.
func UnsafeConvert[T any](origin interface{}) T {
	var dest T

	var bytes []byte
	switch origin := origin.(type) {
	case []byte:
		bytes = origin
	case string:
		bytes = []byte(origin)
	default:
		mbytes, err := json.Marshal(origin)
		if err != nil {
			Error(err)
			return dest
		}
		bytes = mbytes
	}

	err := json.Unmarshal(bytes, &dest)
	if err != nil {
		Error(err)
		return dest
	}

	return dest
}
