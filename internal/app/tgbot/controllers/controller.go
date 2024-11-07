package controllers

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/githubclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgcommons"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/ustorage"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/usecases"
)

// Controller служит для передачи данных в контроллеры
type Controller struct {
	Bot         *tgcommons.Bot
	ModelClient mclient.ModelClient
	TrsUseCases *usecases.TrsUseCases
	Storage     ustorage.UserDataStorage
	GihubClient *githubclient.Client
}

func New(bot *tgcommons.Bot, userStorage ustorage.UserDataStorage, ghToken string) (*Controller, error) {
	mclient, err := mclient.NewMistralClient()
	if err != nil {
		return nil, err
	}

	uc, err := usecases.New()
	if err != nil {
		return nil, err
	}

	ghClient := githubclient.New(ghToken)

	return &Controller{
		Bot:         bot,
		ModelClient: mclient,
		TrsUseCases: uc,
		Storage:     userStorage,
		GihubClient: ghClient,
	}, nil
}
