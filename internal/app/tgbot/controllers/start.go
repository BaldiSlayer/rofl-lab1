package controllers

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
)

// EmptyState - начальное состояние
func (controller *Controller) EmptyState(update tgbotapi.Update) (models.UserState, error) {
	chatID := update.Message.Chat.ID

	err := controller.Bot.SendMessage(chatID, "Введите запрос к базе знаний")
	if err != nil {
		return models.EmptyState, err
	}

	return models.WaitForRequest, err
}

func (controller *Controller) WaitForRequest(update tgbotapi.Update) (models.UserState, error) {
	if update.Message.IsCommand() && update.Message.Command() == "trs" {
		slog.Info("Got command trs", "text", update.Message.Text)

		// TODO: trs handler

		return models.EmptyState, nil
	}

	return controller.handleKnowledgeBaseRequest(update)
}
