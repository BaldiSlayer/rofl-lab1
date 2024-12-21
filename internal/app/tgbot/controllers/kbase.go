package controllers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"

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

func (controller *Controller) getAnswerForKnowledgeBaseRequest(ctx context.Context, userRequest string) (string, error) {
	askResults, err := usecases.AskKnowledgeBase(ctx, controller.ModelClient, userRequest)
	if err != nil {
		return "", fmt.Errorf("failed to ask knowledge base: %w", err)
	}

	gistLink, err := usecases.UploadKnowledgeBaseAnswers(ctx, controller.GihubClient, userRequest, askResults)
	if err != nil {
		return "", fmt.Errorf("failed to ask knowledge base: %w", err)
	}

	answer := tgbotapi.EscapeText(tgbotapi.ModeMarkdown, askResults.Answers[0].Answer)
	buildVersion := version.BuildVersionWithLink()

	message := fmt.Sprintf("%s\n\n[ссылка на использованный контекст](%s)\n\n%s", answer, gistLink, buildVersion)

	slog.Debug(message)

	return message, err
}

func (controller *Controller) handleKnowledgeBaseRequest(
	ctx context.Context,
	update tgbotapi.Update,
) (models.UserState, error) {
	msgID, err := controller.Bot.SendMessageWithReturningID(
		update.Message.Chat.ID,
		"Запрос обрабатывается. Ожидайте.",
	)
	if err != nil {
		return 0, fmt.Errorf("error while trying to send message: %w", err)
	}

	message, err := controller.getAnswerForKnowledgeBaseRequest(ctx, update.Message.Text)
	if err != nil {
		return 0, fmt.Errorf("failed to get answer for knowledge base request: %w", err)
	}

	return models.GetRequest, controller.Bot.EditMarkdownMessage(
		update.Message.Chat.ID,
		msgID,
		message,
	)
}
