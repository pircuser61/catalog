package test_config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "QA"

type Config struct {
	Host   string `split_words:"true" default:"localhost:6000"`
	DbHost string `split_words:"true" default:"localhost"`
	DbPort string `split_words:"true" default:"5432"`
	User   string `split_words:"true" default:"test"`
	Passw  string `split_words:"true" default:"test"`
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(envPrefix, cfg)
	return cfg, err
}

func (cfg *Config) DdConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.User, cfg.Passw, "test")
}
