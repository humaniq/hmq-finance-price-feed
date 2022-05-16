package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Feed struct {
	Assets     []string          `yaml:"assets"`
	Currencies map[string]string `yaml:"currencies"`
	Providers  []ProviderConfig  `yaml:"providers"`
	Diffs      DiffChecker       `yaml:"diffs"`
}

type ProviderConfig struct {
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"`
	Symbols    map[string]string `yaml:"symbols"`
	Currencies map[string]string `yaml:"currencies"`
}

func FeedConfigFromFile(filePath string) (*Feed, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Feed
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

type DiffChecker map[string]int
