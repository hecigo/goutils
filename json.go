package goutils

import (
	"errors"
	"fmt"
	"reflect"

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

	if reflect.TypeOf(origin) == reflect.TypeOf(dest) {
		return origin.(T), nil
	}

	var (
		bytes []byte
		err   error
	)
	switch origin := origin.(type) {
	case []byte:
		bytes = origin
	case string:
		bytes = []byte(origin)
	case map[string]string:
		// dest is struct
		if reflect.TypeOf(dest).Kind() == reflect.Struct {
			// convert origin to map[string]interface{}
			origin2 := make(map[string]interface{})
			for k, v := range origin {
				m := make(map[string]interface{})
				err := json.Unmarshal([]byte(v), &m)
				if err != nil {
					origin2[k] = v
				} else {
					origin2[k] = m
				}
			}

			bytes, err = json.Marshal(origin2)
		} else {
			bytes, err = json.Marshal(origin)
		}

		if err != nil {
			return dest, err
		}
	default:
		mbytes, err := json.Marshal(origin)
		if err != nil {
			return dest, err
		}
		bytes = mbytes
	}

	err = json.Unmarshal(bytes, &dest)
	return dest, err
}

func UnmarshalIntf(origin interface{}, dest interface{}) error {
	if dest == nil {
		return errors.New("failed to unmarshal JSON: unknown destination type")
	}

	var bytes []byte
	switch origin := origin.(type) {
	case []byte:
		bytes = origin
	case string:
		bytes = []byte(origin)
	case map[string]interface{}:
		mbytes, err := json.Marshal(origin)
		if err != nil {
			return err
		}
		bytes = mbytes
	case map[string]string:
		mbytes, err := json.Marshal(origin)
		if err != nil {
			return err
		}
		bytes = mbytes
	case []map[string]interface{}:
		mbytes, err := json.Marshal(origin)
		if err != nil {
			return err
		}
		bytes = mbytes
	default:
		return fmt.Errorf("failed to unmarshal JSON: origin type is invalid\n%v", origin)
	}

	err := json.Unmarshal(bytes, &dest)
	if err != nil {
		mbytes, e := json.Marshal(origin)
		if e != nil {
			return fmt.Errorf("failed to unmarshal JSON: %v\n%v", e, origin)
		}
		return fmt.Errorf("failed to unmarshal JSON: %v\n%v", err, string(mbytes))
	}

	return nil
}
