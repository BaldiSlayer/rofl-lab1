package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
)

// EmptyState - начальное состояние
func (controller *Controller) Start(update tgbotapi.Update) (models.UserState, error) {
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

func (controller *Controller) GetRequest(update tgbotapi.Update) (models.UserState, error) {
	if update.Message == nil {
		return models.GetRequest, nil
	}

	if update.Message.IsCommand() && update.Message.Command() == "trs" {
		return controller.handleTrsRequest(update)
	}

	return controller.handleKnowledgeBaseRequest(update)
}
