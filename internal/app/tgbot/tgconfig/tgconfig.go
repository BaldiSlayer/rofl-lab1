package tgconfig

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type TGBotConfig struct {
	Token string `env:"TGBOTOKEN"`
}

func LoadTGBotConfig() (*TGBotConfig, error) {
	cfg := &TGBotConfig{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env tgbotconfig: %w", err)
	}

	return cfg, err
}
