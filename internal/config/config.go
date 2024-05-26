package config

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-defaults"
)

var (
	VERSION = "dev"
	COMMIT  = "unknown"
)

var cfg Config

type Config struct {
	PostgresConfig PostgresConfig `envconfig:"DATABASE"`
	RedisConfig    RedisConfig    `envconfig:"REDIS"`
}

func (cfg Config) String() string {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (cfg Config) Validate() error {
	if err := cfg.PostgresConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func Init() (*Config, error) {
	defaults.SetDefaults(&cfg)

	err := envconfig.Process("WEBHOOKX", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
