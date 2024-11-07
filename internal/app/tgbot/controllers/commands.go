package controllers

import (
	"context"
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/version"
)

const helpMessage = `Для запроса к Базе Знаний введите запрос

Для проверки завершимости TRS введите:

`+"```"+`
/trs [описание TRS]
`+"```"+`

Пример:

`+"```"+`
/trs
Дана система переписывания термов (TRS): f(x)=a, g(x)=f(f(x)), u(x,y)=c(g(x),f(y)). Я интерпретирую её конструкторы так: a=1, f(x)=x2+2*x+1, g(x)=x3, u(x,y)=x+y, c(x,y)=x+y. Доказывает ли моя интерпретация завершимость trs?
`+"```"

func (controller *Controller) StartCommand(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	return controller.Start(ctx, update)
}

func (controller *Controller) HelpCommand(_ context.Context, update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID
	return models.GetRequest, errors.Join(
		controller.Bot.SendMessage(userID, helpMessage),
		controller.Bot.SendMessage(userID, "Введите запрос к Базе Знаний"),
	)
}

func (controller *Controller) TrsCommand(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		return models.GetTrs, controller.Bot.SendMessage(userID, "Введите TRS")
	}

	return controller.extractTrs(ctx, args, update)
}

func (controller *Controller) VersionCommand(_ context.Context, update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID
	return models.GetRequest, errors.Join(
		controller.Bot.SendMessage(userID, version.BuildVersion()),
		controller.Bot.SendMessage(userID, "Введите запрос к Базе Знаний"),
	)
}
