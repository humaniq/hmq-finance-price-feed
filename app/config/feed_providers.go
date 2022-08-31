package config

import "time"

type PriceProvider struct {
	Name        string               `yaml:"name"`
	EveryString string               `yaml:"every"`
	PancakeSwap *PancakeSwapProvider `yaml:"pancakeswap,omitempty"`
	CoinGecko   *CoinGeckoProvider   `yaml:"coingecko,omitempty"`
	GeoCurrency *GeoCurrencyProvider `yaml:"geocurrency,omitempty"`
}

func (pp *PriceProvider) Every() time.Duration {
	duration, err := time.ParseDuration(pp.EveryString)
	if err != nil {
		return DefaultPriceFeedTickDuration
	}
	return duration
}

type PancakeSwapProvider struct {
	Symbols     []string `yaml:"symbols"`
	AssetMapper string   `yaml:"asset_mapper"`
}
type CoinGeckoProvider struct {
	Currencies []string `yaml:"currencies"`
	Symbols    []string `yaml:"symbols"`
}
type GeoCurrencyProvider struct {
	ApiKey     string            `yaml:"api_key"`
	Symbols    []string          `yaml:"symbols"`
	Currencies map[string]string `yaml:"currencies"`
}
