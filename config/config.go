package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	GrpcAddr = ":8082"
	HttpAddr = "localhost:8087"

	Topic_create = "good_create"
	Topic_update = "good_update"
	Topic_delete = "good_delete"
	Topic_error  = "errors"

	JaegerHostPort = "localhost:6831"
)

var (
	KafkaBrokers = []string{"localhost:19091", "localhost:29091", "localhost:39091"}
)

type TgbotConfig struct {
	ApiKey string
	Debug  bool
}

type PostgreConfig struct {
	host   string
	port   string
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

func GetPostgresConfig() *PostgreConfig {
	cfg := &PostgreConfig{}
	cfg.dbname = getEnv("POSTGRES_DB", "postgres")
	cfg.host = getEnv("POSTGRES_HOST", "localhost")
	cfg.port = getEnv("POSTGRES_PORT", "5433")
	cfg.user = getEnv("POSTGRES_USER", "user")
	cfg.passw = getEnv("POSTGRES_PASSWORD", "1234")
	return cfg
}

func getEnv(name, defaultValue string) string {
	key, ok := os.LookupEnv(name)
	if !ok || key == "" {
		fmt.Printf("%s not set, using default value.\n", name)
		return defaultValue
	}
	return key
}

func GetConnectionString() string {
	cfg := GetPostgresConfig()
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.host, cfg.port, cfg.user, cfg.passw, cfg.dbname)
}
