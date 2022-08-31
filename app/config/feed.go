package config

import (
	"time"

	"github.com/humaniq/hmq-finance-price-feed/pkg/application"
)

const DefaultPriceFeedTickDuration = time.Minute * 5

type PriceFeed struct {
	Assets     []string          `yaml:"assets"`
	Currencies map[string]string `yaml:"currencies"`
	Providers  []PriceProvider   `yaml:"providers"`
	Consumers  []PriceConsumer   `yaml:"consumers"`
	Thresholds Thresholds        `yaml:"thresholds"`
	AssetFiles AssetFiles        `yaml:"asset_files"`
	AssetsData *AssetsData       `yaml:"-"`
}

type Thresholds struct {
	Default Threshold   `yaml:"default,omitempty"`
	Custom  []Threshold `yaml:"custom,omitempty"`
}
type Threshold struct {
	Symbol           string  `yaml:"symbol,omitempty"`
	Currency         string  `yaml:"currency,omitempty"`
	PercentThreshold float64 `yaml:"percent_threshold,omitempty"`
	TimeThreshold    string  `yaml:"time_threshold,omitempty"`
}

type AssetFiles struct {
	CoinGecko []string `yaml:"coingecko"`
	Eth       []string `yaml:"eth"`
}
type AssetsData struct {
	CoinGecko   *CoinGeckoAssets       `yaml:"-"`
	EthNetworks map[string]*EthNetwork `yaml:"-"`
}

func PriceFeedConfigFromFile(filePath string) (*PriceFeed, error) {
	var config PriceFeed
	if err := application.ReadFromYamlFile(filePath, &config); err != nil {
		return nil, err
	}
	config.OverridesFromEnv()
	assets, err := config.ReadAssetsData()
	if err != nil {
		return nil, err
	}
	config.AssetsData = assets
	return &config, nil
}

func (pf *PriceFeed) OverridesFromEnv() {
	for _, consumer := range pf.Consumers {
		if consumer.PriceOracle != nil {
			consumer.PriceOracle.ClientPrivateKey = application.
				StringsValueEnvOverride(consumer.PriceOracle.ClientPrivateKey, "")
		}
	}
	for _, provider := range pf.Providers {
		if provider.GeoCurrency != nil {
			provider.GeoCurrency.ApiKey = application.
				StringsValueEnvOverride(provider.GeoCurrency.ApiKey, "")
		}
	}
}

func (pf *PriceFeed) ReadAssetsData() (*AssetsData, error) {
	coingecko := NewCoinGeckoAssets()
	for _, file := range pf.AssetFiles.CoinGecko {
		assets, err := AssetsFromFile(file)
		if err != nil {
			return nil, err
		}
		for key, val := range assets.Symbols {
			coingecko.Symbols[key] = val
		}
		for key, val := range assets.Currencies {
			coingecko.Currencies[key] = val
		}
	}
	networks := make(map[string]*EthNetwork)
	for _, file := range pf.AssetFiles.Eth {
		network, err := EthNetworkFromFile(file)
		if err != nil {
			return nil, err
		}
		networks[network.Key] = network
	}
	return &AssetsData{
		CoinGecko:   coingecko,
		EthNetworks: networks,
	}, nil
}
