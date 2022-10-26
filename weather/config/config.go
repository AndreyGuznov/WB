package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"weather/tool/log"

	"github.com/joho/godotenv"
)

var envPaths = []string{".env", "/.env", "/root/.env", "~/.env"}

func init() {
	// Try to load files sequentially, until first found
	for _, fname := range envPaths {
		if err := godotenv.Load(fname); err == nil {
			return
		}
	}

	log.Debug(fmt.Sprintf(".env files not found (lookup paths: %s)", strings.Join(envPaths, ", ")))
}

// Get config value as a string
// Value is taken in following order:
// 1. Environmental variable of the same name
// 2. .env file
// 3. Default value
func Get(name string, def ...string) string {
	val, _ := os.LookupEnv(name)

	if val != "" {
		return val
	}

	if len(def) == 0 {
		panic(fmt.Sprintf("Environment variable %s is not set", name))
	}

	return def[0]
}

// GetOptional returns config value as string, if exists
func GetOptional(name string) (val string, ok bool) {
	val, _ = os.LookupEnv(name)

	if val != "" {
		return val, true
	}

	return val, false
}

// GetInt returns config value as integer
func GetInt(name string, def ...int) int {
	val := Get(name, "")

	if val == "" {
		if len(def) == 0 {
			panic(fmt.Sprintf("Environment variable %s is not set", name))
		}

		return def[0]
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("Environment variable %s: %s is not int: %v", name, val, err))
	}

	return res
}

// GetBool returns config value as boolean
func GetBool(name string, def ...bool) bool {
	val := Get(name, "")

	if val == "" {
		if len(def) == 0 {
			panic(fmt.Sprintf("Environment variable %s is not set", name))
		}

		return def[0]
	}

	if val == "0" || val == "false" {
		return false
	}

	return true
}
