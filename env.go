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
	env = flag.String("env", "", "Environment profile")
)

type EnvDuration time.Duration

// Load environment variables from .env file
func LoadEnv() {
	// Load .env (view more: https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use)
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

// Get application name
func AppName() string {
	return Env("APP_NAME", "hexaboi")
}

// Get application version
func AppVersion() string {
	return Env("APP_VERSION", "v0.0.0")
}

// Get API root path
func APIRootPath() string {
	return Env("APP_ROOT_PATH", "/"+AppVersion()[0:2])
}
