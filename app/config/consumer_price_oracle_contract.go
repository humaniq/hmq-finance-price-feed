package config

type EthNetwork struct {
	Name    string                              `yaml:"name"`
	RawUrl  string                              `yaml:"raw_url"`
	ChainId int64                               `yaml:"chain_id"`
	Symbols map[string]EthNetworkSymbolContract `yaml:"symbols"`
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

type PriceOracleContractConsumer struct {
	NetworkUid         string                     `yaml:"network_uid"`
	Name               string                     `yaml:"name"`
	ContractAddressHex string                     `yaml:"contract_address_hex"`
	ClientPrivateKey   string                     `yaml:"client_private_key"`
	Tokens             []PriceOracleContractToken `yaml:"tokens"`
}
