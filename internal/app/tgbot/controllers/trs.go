package controllers

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
)

func (controller *Controller) handleTrsRequest(update tgbotapi.Update) (models.UserState, error) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		err := controller.Bot.SendMessage(update.Message.From.ID, "Введите TRS")
		if err != nil {
			return models.WaitForRequest, err
		}

		return models.WaitForTRS, nil
	}

	trs, request, err := controller.TrsUseCases.ExtractFormalTrs(args)

}
