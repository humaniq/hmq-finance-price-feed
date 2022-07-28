package config

type SavyPrice struct {
	Decimals         int     `yaml:"decimals"`
	AddressHex       string  `yaml:"address_hex"`
	PercentThreshold float64 `yaml:"percent"`
}

type Savy struct {
}
