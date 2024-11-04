package controllers

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/usecases"
)

func (controller *Controller) handleKnowledgeBaseRequest(update tgbotapi.Update) (models.UserState, error) {
	answer, err := usecases.AskKnowledgeBase(controller.ModelClient, update.Message.Text)
	if err != nil {
		return 0, err
	}

	err = controller.Bot.SendMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Ответ Базы Знаний: %s", answer),
	)
	return models.GetRequest, errors.Join(err, controller.Bot.SendMessage(
		update.Message.Chat.ID,
		"Введите запрос к Базе Знаний",
	))
}
