package config

import (
	"github.com/humaniq/hmq-finance-price-feed/pkg/application"
)

type EthNetwork struct {
	Key     string                              `yaml:"key"`
	Name    string                              `yaml:"name"`
	RawUrl  string                              `yaml:"raw_url"`
	ChainId int64                               `yaml:"chain_id"`
	Symbols map[string]EthNetworkSymbolContract `yaml:"symbols"`
}

func EthNetworkFromFile(filePath string) (*EthNetwork, error) {
	var network EthNetwork
	if err := application.ReadFromYamlFile(filePath, &network); err != nil {
		return nil, err
	}
	return &network, nil
}

type EthNetworkSymbolContract struct {
	AddressHex string  `yaml:"address_hex"`
	Decimals   float64 `yaml:"decimals"`
}

type PriceOracleContractToken struct {
	Symbol           string  `yaml:"symbol"`
	Currency         string  `yaml:"currency"`
	PercentThreshold float64 `yaml:"percent_threshold"`
}
