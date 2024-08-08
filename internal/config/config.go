package config

import (
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DbAddr		string	`mapstructure:"DB_ADDR"`
	ListenAddr	string	`mapstructure:"LISTEN_ADDR"`
	Db			string	`mapstructure:"DB"`
	DbUser		string	`mapstructure:"DB_USER"`
	DbPass		string	`mapstructure:"DB_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	log.Info("Loading config from environment")
	var config Config

	viper.AutomaticEnv()

	log.Debug("Parsing environment variables to config struct")
	if err := viper.Unmarshal(&config); err != nil {
		log.Errorf("Error parsing environment variables to config struct: %s", err.Error())
		return nil, err
	}

	log.Info("Successfully loaded config from environment")
	return &config, nil
}