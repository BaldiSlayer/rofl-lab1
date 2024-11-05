package controllers

import (
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/version"
)

const helpMessage = `Для запроса к Базе Знаний введите запрос

Для проверки завершимости TRS введите:
`+"```"+`
/trs [опциональное описание TRS]
`+ "```"

func (controller *Controller) StartCommand(update tgbotapi.Update) (models.UserState, error) {
	return controller.Start(update)
}

func (controller *Controller) HelpCommand(update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID
	return models.GetRequest, errors.Join(
		controller.Bot.SendMessage(userID, helpMessage),
		controller.Bot.SendMessage(userID, "Введите запрос к Базе Знаний"),
	)
}

func (controller *Controller) TrsCommand(update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		return models.GetTrs, controller.Bot.SendMessage(userID, "Введите TRS")
	}

	return controller.extractTrs(args, update)
}

func (controller *Controller) VersionCommand(update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID
	return models.GetRequest, errors.Join(
		controller.Bot.SendMessage(userID, version.BuildVersion()),
		controller.Bot.SendMessage(userID, "Введите запрос к Базе Знаний"),
	)
}
