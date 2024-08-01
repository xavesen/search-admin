package config

import (
	"github.com/spf13/viper"
)

type DbConfig struct {
	DbAddr string `mapstructure:"DB_ADDR"`
	DbPass string `mapstructure:"DB_PASSWORD"`
	Db     int    `mapstructure:"DB"`
}

type Config struct {
	Db         DbConfig
	ListenAddr string `mapstructure:"LISTEN_ADDR"`
}

func LoadConfig() (*Config, error) {
	var config Config

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
