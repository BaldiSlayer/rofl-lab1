package controllers

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgcommons"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/ustorage"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/usecases"
)

// Controller служит для передачи данных в контроллеры
type Controller struct {
	Bot         *tgcommons.Bot
	ModelClient mclient.ModelClient
	TrsUseCases usecases.TrsUseCases
	Storage ustorage.UserDataStorage
}
