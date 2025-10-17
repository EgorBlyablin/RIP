package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Service struct {
		Host string
		Port int
	}

	S3 struct {
		Host string
		Port int
	}
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("failed to read config: %v", err)
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Errorf("failed to unmarshal config: %v", err)
		return nil, err
	}

	log.Info("config parsed")

	return cfg, nil
}
