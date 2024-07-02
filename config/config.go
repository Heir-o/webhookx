package config

import (
	"encoding/json"
	"github.com/creasty/defaults"
	"github.com/kelseyhightower/envconfig"
)

var (
	VERSION = "dev"
	COMMIT  = "unknown"
)

var cfg Config

type Config struct {
	Log            LogConfig      `envconfig:"LOG"`
	PostgresConfig PostgresConfig `envconfig:"DATABASE"`
	RedisConfig    RedisConfig    `envconfig:"REDIS"`
	ServerConfig   ServerConfig   `envconfig:"SERVER"`
}

func (cfg Config) String() string {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (cfg Config) Validate() error {
	if err := cfg.Log.Validate(); err != nil {
		return err
	}
	if err := cfg.PostgresConfig.Validate(); err != nil {
		return err
	}
	if err := cfg.RedisConfig.Validate(); err != nil {
		return err
	}
	if err := cfg.ServerConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func Init() (*Config, error) {
	if err := defaults.Set(&cfg); err != nil {
		return nil, err
	}

	err := envconfig.Process("WEBHOOKX", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
