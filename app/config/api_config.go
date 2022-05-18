package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Api struct {
	Assets     []string          `yaml:"assets"`
	Currencies map[string]string `yaml:"currencies"`
}

func ApiConfigFromFile(filePath string) (*Api, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Api
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
