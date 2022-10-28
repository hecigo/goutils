// `goutils` is a Go utility library that provides a set of useful functions and tools for Go developers.
//
// Including:
//   - `env`: Environment variables management
//   - `log`: Logging
//   - `http`: API Response and Error
//   - `json`: JSON Marshal and Unmarshal everything. `UnsafeConvert` is a function that converts any type to any type.
//   - `time`: utility functions for time, like Now(), Today(), Yesterday()... in UTC+7
//   - `uuid`: UUID generator
//   - `string`: String utility functions: remove Vietnamese accents, format string as URL...
//   - `slice`: De-duplicate, remove empty elements, shuffle, sort...
//   - `struct`: Access struct fields by name
package goutils

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	env = flag.String("env", "", "Environment profile") // Read environment profile from command line
)

// Load environment variables from .env file.
// If the environment variable is already set, it WILL NOT be overwritten.
// The following table shows the priority of environment variables:
// 1. `.env.development.local`
// 2. `.env.test.local`
// 3. `.env.production.local`
// 4. `.env.local`
// 5. `.env.development`
// 6. `.env.test`
// 7. `.env.production`
// 8. `.env`
// View more: https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use
func LoadEnv() {
	if env == nil || *env == "" {
		*env = "development"
	}
	godotenv.Load(".env." + *env + ".local")
	if *env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + *env)
	godotenv.Load() // Load the default environment
}

// Get environment variable. If the environment variable is not set, return the default value.
// T is the type of the environment variable value, and the default value must be the same type.
// T can be string, int, bool, time.Duration, []string or []int
func Env[T string | bool | int | time.Duration | []string | []int](key string, fallback T) T {
	value, ok := os.LookupEnv(strings.ToUpper(key))
	if !ok {
		return fallback
	}

	typeof := reflect.TypeOf(fallback)
	switch typeof.Kind() {
	case reflect.String:
		return reflect.ValueOf(value).Convert(typeof).Interface().(T)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return reflect.ValueOf(val).Convert(typeof).Interface().(T)
	case reflect.Int:
		val, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return reflect.ValueOf(val).Convert(typeof).Interface().(T)
	case reflect.Slice:
		switch typeof.Elem().Kind() {
		case reflect.String:
			return reflect.ValueOf(strings.Split(value, ",")).Convert(typeof).Interface().(T)
		case reflect.Int:
			vals := strings.Split(value, ",")
			var intVals []int
			for _, v := range vals {
				val, err := strconv.Atoi(v)
				if err != nil {
					return fallback
				}
				intVals = append(intVals, val)
			}
			return reflect.ValueOf(intVals).Convert(typeof).Interface().(T)
		}
		return fallback
	case reflect.Struct:
		if typeof == reflect.TypeOf(time.Duration(0)) {
			val, err := time.ParseDuration(value)
			if err != nil {
				return fallback
			}
			return reflect.ValueOf(val).Convert(typeof).Interface().(T)
		}
		return fallback
	default:
		return fallback
	}
}

// Get application name, default is "app"
func AppName() string {
	return Env("APP_NAME", "app")
}

// Get application version, default is "v0.0.0"
func AppVersion() string {
	return Env("APP_VERSION", "v0.0.0")
}

// Get API root path, default is "/v0"
func APIRootPath() string {
	return Env("APP_ROOT_PATH", "/"+AppVersion()[0:2])
}
