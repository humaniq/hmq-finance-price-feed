package config

import "github.com/humaniq/hmq-finance-price-feed/pkg/application"

type CoinGeckoAssets struct {
	Currencies map[string]string `yaml:"currencies"`
	Symbols    map[string]string `yaml:"symbols"`
}

func NewCoinGeckoAssets() *CoinGeckoAssets {
	return &CoinGeckoAssets{
		Currencies: make(map[string]string),
		Symbols:    make(map[string]string),
	}
}
func AssetsFromFile(path string) (*CoinGeckoAssets, error) {
	var config CoinGeckoAssets
	if err := application.ReadFromYamlFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
