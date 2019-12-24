package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	APPENV       string `envconfig:"APPENV"`
	DatastoreURL string `envconfig:"DATASTORE_URL"`
	Port         int    `envconfig:"SERVICE_PORT"`
}

func init() {
	viper.SetDefault("PORT", "80")
	viper.SetDefault("APPENV", "beta")
}

func LoadEnv() (*Config, error) {
	var instance Config
	if err := viper.Unmarshal(&instance); err != nil {
		return nil, err
	}

	err := envconfig.Process("", &instance)
	if err != nil {
		return &instance, err
	}

	return &instance, nil
}
