package goutils

import (
	"errors"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/goccy/go-json"
)

// Convert string to any type by [reflect.Type].
func ReflectStrConv(val string, reflectType reflect.Type) (r interface{}, err error) {
	switch reflectType {
	case reflect.TypeOf(time.Time{}):
		r, err = ParseTime(val)
	case reflect.TypeOf(time.Duration(0)):
		r, err = time.ParseDuration(val)
	default:
		switch reflectType.Kind() {
		case reflect.String:
			r = val
		case reflect.Int:
			r, err = strconv.Atoi(val)
		case reflect.Int8:
			r, err = strconv.ParseInt(val, 10, 8)
		case reflect.Int16:
			r, err = strconv.ParseInt(val, 10, 16)
		case reflect.Int32:
			r, err = strconv.ParseInt(val, 10, 32)
		case reflect.Int64:
			r, err = strconv.ParseInt(val, 10, 64)
		case reflect.Uint:
			r, err = strconv.ParseUint(val, 10, 0)
		case reflect.Uint8:
			r, err = strconv.ParseUint(val, 10, 8)
		case reflect.Uint16:
			r, err = strconv.ParseUint(val, 10, 16)
		case reflect.Uint32:
			r, err = strconv.ParseUint(val, 10, 32)
		case reflect.Uint64:
			r, err = strconv.ParseUint(val, 10, 64)
		case reflect.Float32:
			r, err = strconv.ParseFloat(val, 32)
		case reflect.Float64:
			r, err = strconv.ParseFloat(val, 64)
		case reflect.Bool:
			r, err = strconv.ParseBool(val)
		case reflect.Slice:
			r = reflect.MakeSlice(reflectType, 0, 0).Interface()
			err = json.Unmarshal([]byte(val), &r)
		case reflect.Map:
			r = reflect.MakeMap(reflectType).Interface()
			err = json.Unmarshal([]byte(val), &r)
		case reflect.Struct:
			r = reflect.New(reflectType).Interface()
			err = json.Unmarshal([]byte(val), &r)
		default:
			err = errors.New("unsupported type")
		}
	}
	return r, err
}

// Convert string to any type
func StrConv[T any](val string) (t T, err error) {
	temp, err := ReflectStrConv(val, reflect.TypeOf(t))
	if err != nil {
		return t, err
	}
	return Unmarshal[T](temp)
}

// Convert []string to []T by [ReflectStrConv].
func ReflectSliceStrConv(s []string, reflectType reflect.Type) (result []interface{}, err error) {
	result = make([]interface{}, len(s))
	for i, v := range s {
		temp, err := ReflectStrConv(v, reflectType)
		if err != nil {
			return nil, err
		}
		result[i] = temp
	}
	return result, nil
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

// Convert []string to T. T is a slice of any type.
// T will be convert by [StrConv].
func SliceStrConv[T any](s []string) (result []interface{}, err error) {
	result = make([]interface{}, len(s))
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
	if reflect.TypeOf(new(T)) == reflect.TypeOf(s) {
		return *(*[]*T)(unsafe.Pointer(&s)), nil
	}

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
