package test_config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "QA"

type Config struct {
	GrpcHost     string `split_words:"true" default:"localhost:6000"`
	DbHost       string `split_words:"true" default:"localhost"`
	DbPort       string `split_words:"true" default:"5434"`
	User         string `split_words:"true" default:"user"`
	Passw        string `split_words:"true" default:"1234"`
	DbName       string `split_words:"true" default:"postgres_test"`
	KafkaBrokers []string
}

func FromEnv() (*Config, error) {
	cfg := &Config{KafkaBrokers: []string{"localhost:9095", "localhost:9096", "localhost:9097"}}
	err := envconfig.Process(envPrefix, cfg)
	return cfg, err
}

func (cfg *Config) DdConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.User, cfg.Passw, cfg.DbName)
}
