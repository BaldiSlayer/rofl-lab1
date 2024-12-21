package controllers

import (
	"context"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"log/slog"
	"strings"

	commons "github.com/BaldiSlayer/rofl-lab1/internal/app/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/usecases"
	"github.com/BaldiSlayer/rofl-lab1/internal/version"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (controller *Controller) GetRequest(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	if update.Message == nil {
		return models.GetRequest, nil
	}

	return controller.handleKnowledgeBaseRequest(ctx, update)
}

func (controller *Controller) getAnswerForKnowledgeBaseRequest(
	ctx context.Context,
	userRequest string,
	requests []commons.ModelRequest,
) (string, error) {
	askResults, err := usecases.AskKnowledgeBase(ctx, controller.ModelClient, userRequest, requests)
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

	message, err := controller.getAnswerForKnowledgeBaseRequest(
		ctx,
		update.Message.Text,
		mclient.GetFastModelRequestsPattern(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get answer for knowledge base request: %w", err)
	}

	return models.GetRequest, controller.Bot.EditMarkdownMessage(
		update.Message.Chat.ID,
		msgID,
		message,
	)
}

func (controller *Controller) GetRequestMultiModels(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID

	args := strings.TrimSpace(update.Message.Text)

	if args == "" {
		return models.GetQuestionMultiModels, controller.Bot.SendMessage(userID, "Введите вопрос к базе знаний")
	}

	msgID, err := controller.Bot.SendMessageWithReturningID(
		update.Message.Chat.ID,
		"Запрос обрабатывается. Ожидайте.",
	)
	if err != nil {
		return 0, fmt.Errorf("error while trying to send message: %w", err)
	}

	message, err := controller.getAnswerForKnowledgeBaseRequest(
		ctx,
		update.Message.Text,
		mclient.GetDefaultModelRequestsPattern(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get answer for knowledge base request: %w", err)
	}

	return models.GetRequest, controller.Bot.EditMarkdownMessage(
		update.Message.Chat.ID,
		msgID,
		message,
	)
}

func (controller *Controller) GetCommandMultiModels(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID

	return models.GetQuestionMultiModels,
		controller.Bot.SendMessage(userID, "Введите вопрос к базе знаний. Для ответа на вопрос "+
			"будет сделано 3 разных запроса.")
}
