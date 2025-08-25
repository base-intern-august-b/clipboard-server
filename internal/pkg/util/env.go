package util

import (
	"errors"
	"os"
)

func GetEnvOrDefault(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return v
}

func GetEnvOrErr(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.New("environment variable " + key + " is not set")
	}

	return v, nil
}
