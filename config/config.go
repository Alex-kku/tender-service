package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		PG   `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Address string `env-required:"true" yaml:"address" env:"SERVER_ADDRESS"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		MaxPoolSize int    `env-required:"true" yaml:"max_pool_size" env:"PG_MAX_POOL_SIZE"`
		URL         string `env-required:"true" yaml:"url"           env:"POSTGRES_CONN"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join(configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return cfg, nil
}
