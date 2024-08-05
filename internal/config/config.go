package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbAddr		string	`mapstructure:"DB_ADDR"`
	ListenAddr	string	`mapstructure:"LISTEN_ADDR"`
	Db			string	`mapstructure:"DB"`
	DbUser		string	`mapstructure:"DB_USER"`
	DbPass		string	`mapstructure:"DB_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	var config Config

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}