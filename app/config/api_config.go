package config

import (
	"github.com/humaniq/hmq-finance-price-feed/pkg/application"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
)

type Api struct {
	API        httpapi.Config    `yaml:"api"`
	Assets     []string          `yaml:"assets"`
	Currencies map[string]string `yaml:"currencies"`
	Backend    Backend           `yaml:"backend"`
}

func ApiConfigFromFile(filePath string) (*Api, error) {
	var config Api
	if err := application.ReadFromYamlFile(filePath, &config); err != nil {
		return nil, err
	}
	config.OverridesFromEnv()
	return &config, nil
}

func (a *Api) OverridesFromEnv() {
	overrideAPIFromEnv(&a.API)
	overrideGDSFromEnv(&a.Backend.GoogleDataStore)
}
func overrideAPIFromEnv(a *httpapi.Config) {
	a.ListenString = application.StringsValueEnvOverride(a.ListenString, "LISTEN")
	a.PortString = application.StringsValueEnvOverride(a.PortString, "PORT")
	a.OpenapiPath = application.StringsValueEnvOverride(a.OpenapiPath, "OPENAPI_PATH")
	a.BaseUrl = application.StringsValueEnvOverride(a.BaseUrl, "BASE_URL")
}

type Backend struct {
	GoogleDataStore GoogleDataStore `yaml:"gcloud_datastore"`
}
