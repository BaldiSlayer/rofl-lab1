package controllers

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/usecases"
	"github.com/BaldiSlayer/rofl-lab1/internal/version"
)

func (controller *Controller) GetRequest(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	if update.Message == nil {
		return models.GetRequest, nil
	}

	return controller.handleKnowledgeBaseRequest(ctx, update)
}

func (controller *Controller) handleKnowledgeBaseRequest(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	userRequest := update.Message.Text

	answers, err := usecases.AskKnowledgeBase(ctx, controller.ModelClient, userRequest)
	if err != nil {
		return 0, err
	}

	gistLink, err := usecases.UploadKnowledgeBaseAnswers(ctx, controller.GihubClient, userRequest, answers)
	if err != nil {
		return 0, err
	}

	return models.GetRequest, errors.Join(
		controller.Bot.SendMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("%s\n\n[ссылка на использованный контекст](%s)\n\n%s", answers[0].Answer, gistLink, version.BuildVersion()),
		),
		controller.Bot.SendMessage(
			update.Message.Chat.ID,
			"Введите запрос к Базе Знаний",
		),
	)
}
