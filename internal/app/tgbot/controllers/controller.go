package controllers

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgcommons"
)

// Controller служит для передачи данных в контроллеры
type Controller struct {
	Bot         *tgcommons.Bot
	ModelClient mclient.ModelClient
}
