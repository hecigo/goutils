package goutils

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// Convert string to any type
func StrConv[T any](val string) (t T, err error) {
	tType := reflect.TypeOf(t)

	var temp interface{}
	switch tType {
	case reflect.TypeOf(time.Time{}):
		temp, err = ParseTime(val)
	case reflect.TypeOf(time.Duration(0)):
		temp, err = time.ParseDuration(val)
	default:
		switch tType.Kind() {
		case reflect.Int:
			temp, err = strconv.Atoi(val)
		case reflect.Int8:
			temp, err = strconv.ParseInt(val, 10, 8)
		case reflect.Int16:
			temp, err = strconv.ParseInt(val, 10, 16)
		case reflect.Int32:
			temp, err = strconv.ParseInt(val, 10, 32)
		case reflect.Int64:
			temp, err = strconv.ParseInt(val, 10, 64)
		case reflect.Uint:
			temp, err = strconv.ParseUint(val, 10, 0)
		case reflect.Uint8:
			temp, err = strconv.ParseUint(val, 10, 8)
		case reflect.Uint16:
			temp, err = strconv.ParseUint(val, 10, 16)
		case reflect.Uint32:
			temp, err = strconv.ParseUint(val, 10, 32)
		case reflect.Uint64:
			temp, err = strconv.ParseUint(val, 10, 64)
		case reflect.Float32:
			temp, err = strconv.ParseFloat(val, 32)
		case reflect.Float64:
			temp, err = strconv.ParseFloat(val, 64)
		case reflect.Bool:
			temp, err = strconv.ParseBool(val)
		case reflect.String:
			temp = val
		case reflect.Struct, reflect.Slice, reflect.Map:
			return Unmarshal[T](val)
		default:
			return t, errors.New("unsupported type")
		}
	}
	return temp.(T), err
}

// Convert map[string]string to map[string]T. T can be any type.
// T will be convert by [StrConv].
func MapStrConv[T any](m map[string]string) (result map[string]T, err error) {
	result = make(map[string]T)
	for k, v := range m {
		t, err := StrConv[T](v)
		if err != nil {
			return nil, err
		}
		result[k] = t
	}
	return result, nil
}

// Convert map[string]*string to map[string]*T. T can be any type.
// T will be convert by [StrConv].
func MapPtrStrConv[T any](m map[string]*string) (result map[string]*T, err error) {
	result = make(map[string]*T)
	for k, v := range m {
		if v == nil {
			result[k] = nil
			continue
		}

		t, err := StrConv[T](*v)
		if err != nil {
			return nil, err
		}
		result[k] = &t
	}
	return result, nil
}

// Convert []string to []T. T can be any type.
// T will be convert by [StrConv].
func SliceStrConv[T any](s []string) (result []T, err error) {
	result = make([]T, len(s))
	for i, v := range s {
		t, err := StrConv[T](v)
		if err != nil {
			return nil, err
		}
		result[i] = t
	}
	return result, nil
}

// Convert []*string to []*T. T can be any type.
// T will be convert by [StrConv].
func SlicePtrStrConv[T any](s []*string) (result []*T, err error) {
	result = make([]*T, len(s))
	for i, v := range s {
		if v == nil {
			result[i] = nil
			continue
		}

		t, err := StrConv[T](*v)
		if err != nil {
			return nil, err
		}
		result[i] = &t
	}
	return result, nil
}
