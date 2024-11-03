package controllers

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
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
}

func New(bot *tgcommons.Bot, userStorage ustorage.UserDataStorage) (*Controller, error) {
	context, err := models.LoadQABase()
	if err != nil {
		return nil, err
	}

	mclient, err := mclient.NewMistralClient(context)
	if err != nil {
		return nil, err
	}

	uc, err := usecases.New()
	if err != nil {
		return nil, err
	}

	return &Controller{
		Bot:         bot,
		ModelClient: mclient,
		TrsUseCases: uc,
		Storage:     userStorage,
	}, nil
}
