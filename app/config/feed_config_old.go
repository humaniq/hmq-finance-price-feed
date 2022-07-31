package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type FeedOld struct {
	Assets             []string          `yaml:"assets"`
	Currencies         map[string]string `yaml:"currencies"`
	Providers          []ProviderConfig  `yaml:"prices"`
	Diffs              Diffs             `yaml:"diffs"`
	ForceUpdateSeconds map[string]int64  `yaml:"force_update_seconds"`
}

type ProviderConfig struct {
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"`
	Symbols    map[string]string `yaml:"symbols"`
	Currencies map[string]string `yaml:"currencies"`
}

func FeedCfgFromFile(filePath string) (*FeedOld, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config FeedOld
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

type Diffs map[string]int

func (dc Diffs) Diff(symbol string) int {
	if val, found := dc[symbol]; found {
		return val
	}
	return 10000
}

type TSDiffs map[string]time.Duration

func NewTSDiffsFromSeconds(seconds map[string]int64) TSDiffs {
	result := make(map[string]time.Duration)
	for key, val := range seconds {
		result[key] = time.Duration(val) * time.Second
	}
	return result
}
func (dc TSDiffs) Diff(symbol string) time.Duration {
	if val, found := dc[symbol]; found {
		return val
	}
	return time.Hour * 24
}
