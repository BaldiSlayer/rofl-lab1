package bconfig

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	// Используется для автоматической загрузки переменных, из файла ".env".
	_ "github.com/joho/godotenv/autoload"
)

type BackendConfig struct {
	Port string `env:"PORT" envDefault:"9000"`
}

func LoadBackendConfig() (*BackendConfig, error) {
	cfg := &BackendConfig{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env bconfig: %w", err)
	}

	return cfg, err
}
