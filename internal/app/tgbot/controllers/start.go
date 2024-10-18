package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
)

// EmptyState - начальное состояние, в котором пользователь еще не имеет своего состояния в системе
func (controller *Controller) EmptyState(update tgbotapi.Update) (models.UserState, error) {
	chatID := update.Message.Chat.ID

	err := controller.Bot.SendStartUpKeyboard(chatID)
	if err != nil {
		return models.EmptyState, err
	}

	return models.StartState, err
}

func (controller *Controller) Start(update tgbotapi.Update) (models.UserState, error) {
	if update.Message.Text == "База знаний" {
		err := controller.Bot.RemoveKeyboard(
			update.Message.Chat.ID,
			"Введите свой запрос к LLM:",
		)
		if err != nil {
			return models.EmptyState, err
		}

		return models.WaitForKBQuestion, nil
	}

	if update.Message.Text == "TRS Solver" {
		err := controller.Bot.RemoveKeyboard(
			update.Message.Chat.ID,
			"Введите TRS:",
		)
		if err != nil {
			return models.EmptyState, err
		}

		return models.TRSState, nil
	}

	return models.EmptyState, nil
}
