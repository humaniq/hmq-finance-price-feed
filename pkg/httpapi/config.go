package httpapi

import "fmt"

type Config struct {
	ListenString string `yaml:"listen" env:"LISTEN"`
	PortString   string `yaml:"port" env:"PORT"`
	OpenapiPath  string `yaml:"openapi_path" env:"OPENAPI_PATH"`
	BaseUrl      string `yaml:"base_url" env:"BASE_URL"`
}

func (c *Config) Listen() string {
	return fmt.Sprintf("%s:%s", c.ListenString, c.PortString)
}
