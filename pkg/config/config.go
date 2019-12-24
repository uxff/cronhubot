package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	APPENV       string `envconfig:"APPENV"`

	// like: mysql://yourusername:yourpwd@tcp(yourmysqlhost)/yourdbname?charset=utf8mb4&parseTime=True&loc=Local
	DatastoreURL string `envconfig:"DATASTORE_URL"`

	Port         int    `envconfig:"SERVICE_PORT"`
}

func init() {
	viper.SetDefault("PORT", "80")
	viper.SetDefault("APPENV", "dev")
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
