package controllers

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
)

// EmptyState - начальное состояние
func (controller *Controller) Start(_ context.Context, update tgbotapi.Update) (models.UserState, error) {
	if update.Message == nil {
		return models.Start, nil
	}

	chatID := update.Message.Chat.ID

	err := controller.Bot.SendMessage(chatID, "Введите запрос к Базе Знаний")
	if err != nil {
		return models.Start, err
	}

	return models.GetRequest, err
}
