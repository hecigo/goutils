package goutils

import (
	"reflect"
	"time"

	"github.com/goccy/go-json"
)

type JSON json.RawMessage

// Marshal an interface to a string.
func Marshal(origin interface{}) (string, error) {
	bytes, err := json.Marshal(origin)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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
	case time.Time, time.Duration, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return reflect.ValueOf(origin).Interface().(T), nil
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
