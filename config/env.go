package config

import (
	"fmt"
	"os"
	"strconv"
)

func getEnv(variableName string) (string, error) {
	value := os.Getenv(variableName)
	if value == "" {
		return value, fmt.Errorf("got empty value for '%s' from env", variableName)
	}
	return value, nil
}

func getEnvInt(variableName string) (int, error) {
	value := os.Getenv(variableName)
	result, err := strconv.Atoi(value)
	if err != nil {
		return result, fmt.Errorf("can't convert '%s' env value '%s' to int: %w", variableName, value, err)
	}
	return result, nil
}