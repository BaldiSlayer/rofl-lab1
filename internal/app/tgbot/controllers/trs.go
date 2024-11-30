package controllers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

const (
	confirmCallbackData = "Confirm"
	fixCallbackData     = "Fix"
)

type extractData struct {
	parseError  string
	userRequest string
	formalTrs   string
}

func (controller *Controller) GetTrs(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	if update.Message == nil {
		return models.GetTrs, nil
	}

	return controller.extractTrs(ctx, update.Message.Text, update)
}

func (controller *Controller) extractTrs(ctx context.Context, userRequest string, update tgbotapi.Update) (models.UserState, error) {
	err := controller.Storage.SetRequest(ctx, update.SentFrom().ID, update.Message.Text)
	if err != nil {
		return 0, err
	}

	extractData, err := controller.TrsUseCases.ExtractFormalTrs(ctx, userRequest)
	return controller.handleExctractResult(ctx, update, extractData.Trs, extractData.FormalizedTrs, err)
}

func (controller *Controller) handleExctractResult(ctx context.Context, update tgbotapi.Update, trs trsparser.Trs,
	formalized string, extractError error) (models.UserState, error) {

	userID := update.SentFrom().ID

	err := controller.Storage.SetFormalTRS(ctx, userID, formalized)
	if err != nil {
		return 0, err
	}

	var parseError *trsparser.ParseError
	if errors.As(extractError, &parseError) {
		err := controller.Storage.SetParseError(ctx, userID, parseError.LlmMessage)
		if err != nil {
			return 0, err
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Исправить", fixCallbackData),
		))
		return models.FixTrs, controller.Bot.SendMarkdownMessageWithKeyboard(
			userID,
			fmt.Sprintf("Ошибка при формализации TRS\nРезультат Formalize:\n```\n%s\n```\nРезультат Parse:\n```\n%s\n```\n\n"+
				"Переформулируйте запрос в новом сообщении, либо запустите процесс автоматического исправления с помощью кнопки под этим сообщением",
				formalized, parseError.LlmMessage),
			keyboard,
		)
	}

	if extractError != nil {
		return 0, extractError
	}

	err = controller.Storage.SetTRS(ctx, userID, trs)
	if err != nil {
		return 0, err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Подтвердить", confirmCallbackData),
	))

	return models.ValidateTrs, controller.Bot.SendMarkdownMessageWithKeyboard(userID,
		fmt.Sprintf("Результат формализации:\n```\n%s\n```\n\n"+
			"Подтвердите его с помощью кнопки под этим сообщением, либо опишите ошибку в новом сообщении", toString(trs)), keyboard)
}

func (controller *Controller) ValidateTrs(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID

	if update.CallbackQuery != nil && update.CallbackQuery.Data == confirmCallbackData {
		trs, err := controller.Storage.GetTRS(ctx, userID)
		if err != nil {
			return 0, err
		}

		res, err := controller.TrsUseCases.InterpretFormalTrs(ctx, trs)
		if err != nil {
			return 0, err
		}

		return models.GetRequest, errors.Join(
			controller.Bot.SendMarkdownMessage(userID, fmt.Sprintf("Результат интерпретации TRS:\n```\n%s\n```", res)),
			controller.Bot.SendMessage(userID, "Введите запрос к Базе Знаний"),
		)
	} else if update.Message != nil {
		errorDescription := update.Message.Text

		userRequest, err := controller.Storage.GetRequest(ctx, userID)
		if err != nil {
			return 0, err
		}

		formalTrs, err := controller.Storage.GetFormalTRS(ctx, userID)
		if err != nil {
			return 0, err
		}

		fixData, err := controller.TrsUseCases.FixFormalTrs(ctx, userRequest, formalTrs, errorDescription)
		return controller.handleExctractResult(ctx, update, fixData.Trs, fixData.FormalizedTrs, err)
	}

	return models.ValidateTrs, nil
}

func (controller *Controller) FixTrs(ctx context.Context, update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID

	if update.CallbackQuery != nil && update.CallbackQuery.Data == fixCallbackData {
		extractData, err := controller.getExtractData(ctx, userID)
		if err != nil {
			return 0, err
		}

		fixData, err := controller.TrsUseCases.FixFormalTrs(ctx,
			extractData.userRequest, extractData.formalTrs, extractData.parseError)

		return controller.handleExctractResult(ctx, update, fixData.Trs, fixData.FormalizedTrs, err)
	}

	if update.Message != nil {
		userRequest := update.Message.Text
		return controller.extractTrs(ctx, userRequest, update)
	}

	return models.FixTrs, nil
}

func (controller *Controller) getExtractData(ctx context.Context, userID int64) (extractData, error) {
	parseError, err := controller.Storage.GetParseError(ctx, userID)
	if err != nil {
		return extractData{}, err
	}

	userRequest, err := controller.Storage.GetRequest(ctx, userID)
	if err != nil {
		return extractData{}, err
	}

	formalTrs, err := controller.Storage.GetFormalTRS(ctx, userID)
	if err != nil {
		return extractData{}, err
	}

	return extractData{
		parseError:  parseError,
		userRequest: userRequest,
		formalTrs:   formalTrs,
	}, err
}

func toString(trs trsparser.Trs) string {
	lines := make([]string, 0, len(trs.Interpretations)+len(trs.Rules)+1)

	variables := fmt.Sprintf("variables = %s", strings.Join(trs.Variables, ", "))
	lines = append(lines, variables)

	lines = append(lines, trs.Rules...)
	lines = append(lines, trs.Interpretations...)

	return strings.Join(lines, "\n")
}
