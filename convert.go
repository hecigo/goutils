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
		case reflect.Interface:
			err = errors.New("`reflectType` must be a specific type")
		default:
			err = errors.New("`reflectType` is not supported")
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

// Convert map[string]string to map[string][reflect.Type] by [ReflectStrConv].
func ReflectMapStrConv(m map[string]string, elemType reflect.Type) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	for k, v := range m {
		temp, err := ReflectStrConv(v, elemType)
		if err != nil {
			return nil, err
		}
		result[k] = temp
	}
	return result, nil
}

// Convert map[string]string to T by [ReflectMapStrConv], with T is a map[string]any.
func MapStrConv[T any](m map[string]string) (result T, err error) {
	var t T
	if reflect.TypeOf(m) == reflect.TypeOf(t) {
		return reflect.ValueOf(m).Interface().(T), nil
	}

	temp, err := ReflectMapStrConv(m, reflect.TypeOf(t).Elem())
	if err != nil {
		return result, err
	}
	return Unmarshal[T](temp)
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

// Convert map[string]interface{} to a struct.
func MapStrAnyToStruct(source map[string]interface{}, dest interface{}) error {
	// Convert map to json string
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}

	// Convert json string to struct
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}

	return nil
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

// Convert any to string, throw error
func AnyToStr(val interface{}) (string, error) {
	switch val := val.(type) {
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case int:
		return strconv.Itoa(val), nil
	case int8:
		return strconv.Itoa(int(val)), nil
	case int16:
		return strconv.Itoa(int(val)), nil
	case int32:
		return strconv.Itoa(int(val)), nil
	case int64:
		return strconv.Itoa(int(val)), nil
	case uint:
		return strconv.Itoa(int(val)), nil
	case uint8:
		return strconv.Itoa(int(val)), nil
	case uint16:
		return strconv.Itoa(int(val)), nil
	case uint32:
		return strconv.Itoa(int(val)), nil
	case uint64:
		return strconv.Itoa(int(val)), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(val), nil
	case time.Time:
		return TimeStr(val), nil
	case time.Duration:
		return val.String(), nil
	default:
		return Marshal(val)
	}
}

// Convert any to string with error as a result.
func ToStr(val interface{}) string {
	str, err := AnyToStr(val)
	if err != nil {
		return err.Error()
	}
	return str
}
