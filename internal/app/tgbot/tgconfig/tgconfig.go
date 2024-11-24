package tgconfig

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type TGBotConfig struct {
	TgToken    string `env:"TGTOKEN"`
	GhToken    string `env:"GHTOKEN"`
	PgUser     string `env:"POSTGRES_USER"`
	PgPassword string `env:"POSTGRES_PASSWORD"`
	PgDBName   string `env:"POSTGRES_DB"`
}

func LoadTGBotConfig() (*TGBotConfig, error) {
	cfg := &TGBotConfig{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env tgbotconfig: %w", err)
	}

	return cfg, err
}
