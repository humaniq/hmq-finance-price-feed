package config

import "github.com/humaniq/hmq-finance-price-feed/pkg/application"

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
