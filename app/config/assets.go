package config

import "github.com/humaniq/hmq-finance-price-feed/pkg/application"

type Assets struct {
	Currencies map[string]string `yaml:"currencies"`
	Symbols    map[string]string `yaml:"symbols"`
}

func AssetsFromFile(path string) (*Assets, error) {
	var config Assets
	if err := application.ReadFromYamlFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
