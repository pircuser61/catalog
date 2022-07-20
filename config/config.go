package config

import (
	"errors"
	"os"
)

type TgbotConfig struct {
	ApiKey string
	Debug  bool
}

func GetConfig() (*TgbotConfig, error) {
	var cfg TgbotConfig
	key, ok := os.LookupEnv("CATALOG_TG_API_KEY")
	if !ok || key == "" {
		return nil, errors.New("api key not found, check the CATALOG_TG_API_KEY environment variable")
	}
	cfg.ApiKey = key
	cfg.Debug = true
	return &cfg, nil
}
