package config

import (
	"time"

	"github.com/humaniq/hmq-finance-price-feed/pkg/application"
)

const DefaultPriceFeedTickDuration = time.Minute * 5

type PriceFeed struct {
	Assets      []string          `yaml:"assets"`
	Currencies  map[string]string `yaml:"currencies"`
	Providers   PriceProviders    `yaml:"providers"`
	Consumers   PriceConsumers    `yaml:"consumers"`
	Coingecko   Coingecko         `yaml:"coingecko"`
	EthNetworks EthNetworks       `yaml:"eth"`
}

type EthNetworks struct {
	NetworksPath string                `yaml:"networks_path"`
	Networks     map[string]EthNetwork `yaml:"networks"`
}

func PriceFeedConfigFromFile(filePath string) (*PriceFeed, error) {
	config := PriceFeed{
		Currencies: make(map[string]string),
		EthNetworks: EthNetworks{
			Networks: make(map[string]EthNetwork),
		},
	}
	if err := application.ReadFromYamlFile(filePath, &config); err != nil {
		return nil, err
	}
	if config.EthNetworks.NetworksPath != "" {
		if err := application.ReadFromYamlFile(config.EthNetworks.NetworksPath, &config.EthNetworks.Networks); err != nil {
			return nil, err
		}
	}
	config.OverridesFromEnv()
	return &config, nil
}

func (pf *PriceFeed) OverridesFromEnv() {
	pf.Coingecko.AssetsPath = application.StringsValueEnvOverride(pf.Coingecko.AssetsPath, "COINGECKO_ASSETS_PATH")
	for index, gecko := range pf.Consumers.PriceOracles {
		gecko.ClientPrivateKey = application.StringsValueEnvOverride(gecko.ClientPrivateKey, "")
		pf.Consumers.PriceOracles[index] = gecko
	}
}

type PriceProviders struct {
	Coingeckos []CoingeckoFeedConfig `yaml:"coingecko"`
}

type CoingeckoFeedConfig struct {
	Name        string   `yaml:"name"`
	Currencies  []string `yaml:"currencies"`
	Symbols     []string `yaml:"symbols"`
	EveryString string   `yaml:"every"`
}

func (cfc *CoingeckoFeedConfig) Every() time.Duration {
	duration, err := time.ParseDuration(cfc.EveryString)
	if err != nil {
		return DefaultPriceFeedTickDuration
	}
	return duration
}

type PriceConsumers struct {
	PriceOracles     []PriceOracleContractConsumer `yaml:"price_oracle"`
	StorageConsumers []StorageConsumer             `yaml:"price_storage"`
}

type Coingecko struct {
	AssetsPath string `yaml:"assets_path"`
}
