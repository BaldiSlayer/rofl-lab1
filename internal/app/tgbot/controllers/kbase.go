package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/usecases"
)

const (
	waitForKBQuestionTimeout = 10 * time.Second
)

func (controller *Controller) handleKnowledgeBaseRequest(update tgbotapi.Update) (models.UserState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), waitForKBQuestionTimeout)
	defer cancel()

	// FIXME: антипаттерн так получать значения?
	var answer string
	var err error
	doneChan := make(chan struct{}, 1)
	go func() {
		answer, err = usecases.AskKnowledgeBase(controller.ModelClient, update.Message.Text)
		if err != nil {
			slog.Error(err.Error())
		}
		doneChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		err = controller.Bot.SendMessage(
			update.Message.Chat.ID,
			"Таймаут при запросе к Базе Знаний, введите новый запрос",
		)
		if err != nil {
			return models.WaitForRequest, err
		}
	case <-doneChan:
		if err != nil {
			err = controller.Bot.SendMessage(
				update.Message.Chat.ID,
				"Ошибка при запросе к Базе Знаний, введите новый запрос",
			)
			if err != nil {
				return models.WaitForRequest, err
			}

			return models.WaitForRequest, err
		}

		err = controller.Bot.SendMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("Ответ Базы Знаний: %s", answer),
		)
		if err != nil {
			return models.WaitForRequest, err
		}
	}

	panic("unreachable")
}
