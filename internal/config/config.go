package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type TgbotConfig struct {
	ApiKey string
	Debug  bool
}

type PostgreConfig struct {
	host   string
	port   int
	user   string
	passw  string
	dbname string
}

func GetTgBotConfig() (*TgbotConfig, error) {
	var cfg TgbotConfig
	key, ok := os.LookupEnv("CATALOG_TG_API_KEY")
	if !ok || key == "" {
		return nil, errors.New("api key not found, check the CATALOG_TG_API_KEY environment variable")
	}
	cfg.ApiKey = key
	cfg.Debug = false
	return &cfg, nil
}

func GetPostgreConfig() (*PostgreConfig, error) {
	key, ok := os.LookupEnv("POSTGRE_CONN")
	if !ok || key == "" {
		return nil, errors.New("api key not found, check the POSTGRE_CONN environment variable")
	}
	params := strings.Split(key, ",")
	if len(params) < 2 {
		return nil, errors.New("wrong param in POSTGRE_CONN environment, <login>,<passw>")
	}
	return &PostgreConfig{host: "localhost",
		port: 5432, user: params[0], passw: params[1], dbname: "postgres"}, nil
}

func GetConnectionString() string {
	cfg, err := GetPostgreConfig()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.host, cfg.port, cfg.user, cfg.passw, cfg.dbname)
}
