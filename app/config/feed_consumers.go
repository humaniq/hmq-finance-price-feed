package config

import (
	"github.com/humaniq/hmq-finance-price-feed/pkg/application"
)

type PriceConsumer struct {
	Name            string               `yaml:"name"`
	Tokens          []PriceConsumerToken `yaml:"tokens"`
	PriceOracle     *PriceOracleConsumer `yaml:"price_oracle,omitempty"`
	GoogleDataStore *GoogleDataStore     `yaml:"google_datastore,omitempty"`
}
type PriceConsumerToken struct {
	Symbol   string `yaml:"symbol"`
	Currency string `yaml:"currency"`
}

type PriceOracleConsumer struct {
	NetworkKey         string `yaml:"network_key"`
	ContractAddressHex string `yaml:"contract_address_hex"`
	ClientPrivateKey   string `yaml:"client_private_key"`
}

const DefaultPriceRecordKind = "hmq_price_assets"

type GoogleDataStore struct {
	ProjectIdString             string `yaml:"project_id"`
	PriceAssetsRecordKindString string `yaml:"price_assets_kind"`
}

func (gds *GoogleDataStore) ProjectID() string {
	return gds.ProjectIdString
}
func (gds *GoogleDataStore) PriceAssetsKind() string {
	if gds.PriceAssetsRecordKindString != "" {
		return gds.PriceAssetsRecordKindString
	}
	return DefaultPriceRecordKind
}
func overrideGDSFromEnv(gds *GoogleDataStore) {
	gds.ProjectIdString = application.StringsValueEnvOverride(gds.ProjectIdString, "DATASTORE_PROJECT_ID")
	gds.PriceAssetsRecordKindString = application.StringsValueEnvOverride(gds.PriceAssetsRecordKindString, "DATASTORE_PRICE_ASSETS_KIND")
}

//func (sc *StorageConsumer) TimeDelta() time.Duration {
//	duration, err := time.ParseDuration(sc.TimeDiffStr)
//	if err != nil {
//		return 0
//	}
//	return duration
//}
