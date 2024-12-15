package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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

	askResults, err := usecases.AskKnowledgeBase(ctx, controller.ModelClient, userRequest)
	if err != nil {
		return 0, err
	}

	gistLink, err := usecases.UploadKnowledgeBaseAnswers(ctx, controller.GihubClient, userRequest, askResults)
	if err != nil {
		return 0, err
	}

	answer := tgbotapi.EscapeText(tgbotapi.ModeMarkdown, askResults.Answers[0].Answer)
	buildVersion := version.BuildVersionWithLink()

	message := fmt.Sprintf("%s\n\n[ссылка на использованный контекст](%s)\n\n%s", answer, gistLink, buildVersion)

	slog.Info(message)

	return models.GetRequest, errors.Join(
		controller.Bot.SendMarkdownMessage(
			update.Message.Chat.ID,
			message,
		),
		controller.Bot.SendMessage(
			update.Message.Chat.ID,
			"Введите запрос к Базе Знаний",
		),
	)
}
