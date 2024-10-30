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
		err = errors.Join(err, controller.Bot.SendMessage(
			update.Message.Chat.ID,
			"Ошибка при запросе к Базе Знаний, введите новый запрос",
		))
		return models.WaitForRequest, err
	}

	err = controller.Bot.SendMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Ответ Базы Знаний: %s", answer),
	)
	if err != nil {
		return models.WaitForRequest, err
	}

	return models.WaitForRequest, nil
}
