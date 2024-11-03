package controllers

import (
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

// TODO: отправлять сообщение об ошибке

func (controller *Controller) handleTrsRequest(update tgbotapi.Update) (models.UserState, error) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		return models.GetTrs, controller.Bot.SendMessage(update.Message.From.ID, "Введите TRS")
	}

	return controller.extractTrs(args, update)
}

func (controller *Controller) GetTrs(update tgbotapi.Update) (models.UserState, error) {
	if update.Message == nil {
		return models.GetTrs, nil
	}

	return controller.extractTrs(update.Message.Text, update)
}

func (controller *Controller) extractTrs(userRequest string, update tgbotapi.Update) (models.UserState, error) {
	err := controller.Storage.SetRequest(update.SentFrom().ID, update.Message.Text)
	if err != nil {
		return 0, err
	}

	trs, formalized, err := controller.TrsUseCases.ExtractFormalTrs(userRequest)
	return controller.handleExctractResult(update, trs, formalized, err)
}

func (controller *Controller) handleExctractResult(update tgbotapi.Update, trs trsparser.Trs,
	formalized string, extractError error) (models.UserState, error) {

	userID := update.SentFrom().ID

	err := controller.Storage.SetFormalTRS(userID, formalized)
	if err != nil {
		return 0, err
	}

	var parseError *trsparser.ParseError
	if errors.As(extractError, &parseError) {
		err := controller.Storage.SetParseError(userID, parseError.LlmMessage)
		if err != nil {
			return 0, err
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Исправить", fixCallbackData),
		))
		return models.FixTrs, controller.Bot.SendMessageWithKeyboard(
			userID,
			fmt.Sprintf("Ошибка при формализации TRS\nРезультат Formalize:\n%s\nРезультат Parse:\n%s\n\n"+
				"Переформулируйте запрос в новом сообщении, либо запустите процесс автоматического исправления по кнопке ниже",
				formalized, parseError.LlmMessage),
			keyboard,
		)
	}

	if extractError != nil {
		return 0, extractError
	}

	err = controller.Storage.SetTRS(userID, trs)
	if err != nil {
		return 0, err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Подтвердить", confirmCallbackData),
	))

	return models.ValidateTrs, controller.Bot.SendMessageWithKeyboard(userID,
		fmt.Sprintf("Результат формализации:\n%s\n\n"+
			"Подтвердите его с помощью кнопки ниже, либо опишите ошибку в новом сообщении", toString(trs)), keyboard)
}

func (controller *Controller) ValidateTrs(update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID

	if update.CallbackQuery != nil && update.CallbackQuery.Data == confirmCallbackData {
		trs, err := controller.Storage.GetTRS(userID)
		if err != nil {
			return 0, err
		}

		res, err := controller.TrsUseCases.InterpretFormalTrs(trs)
		if err != nil {
			return 0, err
		}

		_ = controller.Bot.SendMessage(userID, fmt.Sprintf("Результат интерпретации TRS:\n%s", res))

		return models.GetRequest, controller.Bot.SendMessage(userID, "Введите запрос к Базе Знаний")
	} else if update.Message != nil {
		errorDescription := update.Message.Text

		userRequest, err := controller.Storage.GetRequest(userID)
		if err != nil {
			return 0, err
		}

		formalTrs, err := controller.Storage.GetFormalTRS(userID)
		if err != nil {
			return 0, err
		}

		trs, formalTrs, err := controller.TrsUseCases.FixFormalTrs(userRequest, formalTrs, errorDescription)
		return controller.handleExctractResult(update, trs, formalTrs, err)
	}

	return models.ValidateTrs, nil
}

func (controller *Controller) FixTrs(update tgbotapi.Update) (models.UserState, error) {
	userID := update.SentFrom().ID

	if update.CallbackQuery != nil && update.CallbackQuery.Data == fixCallbackData {
		parseError, userRequest, formalTrs, err := controller.getExtractData(userID)
		if err != nil {
			return 0, err
		}

		trs, formalTrs, err := controller.TrsUseCases.FixFormalTrs(userRequest, formalTrs, parseError)

		return controller.handleExctractResult(update, trs, formalTrs, err)
	} else if update.Message != nil {
		userRequest := update.Message.Text

		return controller.extractTrs(userRequest, update)
	}

	return models.FixTrs, nil
}

func (controller *Controller) getExtractData(userID int64) (string, string, string, error) {
	parseError, err := controller.Storage.GetParseError(userID)
	if err != nil {
		return "", "", "", err
	}

	userRequest, err := controller.Storage.GetRequest(userID)
	if err != nil {
		return "", "", "", err
	}

	formalTrs, err := controller.Storage.GetFormalTRS(userID)
	if err != nil {
		return "", "", "", err
	}

	return parseError, userRequest, formalTrs, nil
}

func toString(trs trsparser.Trs) string {
	var lines []string

	variables := fmt.Sprintf("variables = %s", strings.Join(trs.Variables, ", "))
	lines = append(lines, variables)

	lines = append(lines, trs.Rules...)
	lines = append(lines, trs.Interpretations...)

	return strings.Join(lines, "\n")
}
