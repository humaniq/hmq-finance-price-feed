package application

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadFromYamlFile(filePath string, value interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(value); err != nil {
		return err
	}
	return nil
}

func StringsValueEnvOverride(value string, envKey string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	pair := strings.Split(value, "|")
	if len(pair) > 1 {
		if valueFromEnv := os.Getenv(pair[1]); valueFromEnv != "" {
			return valueFromEnv
		}
		return pair[0]
	}
	return value
}
