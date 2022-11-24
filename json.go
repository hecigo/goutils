package goutils

import (
	"github.com/goccy/go-json"
)

type JSON json.RawMessage

// Marshal an interface to a string.
func Marshal(origin interface{}) string {
	bytes, err := json.Marshal(origin)
	if err != nil {
		Error(err)
		return ""
	}
	return string(bytes)
}

// Convert a string or struct to another struct.
// This func uses JSON as a middle data type to convert.
func Unmarshal[T any](origin interface{}) (T, error) {
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
			return dest, err
		}
		bytes = mbytes
	}

	err := json.Unmarshal(bytes, &dest)
	if err != nil {
		return dest, err
	}

	return dest, nil
}
