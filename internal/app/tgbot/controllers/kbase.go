package controllers

import (
	"context"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (controller *Controller) WaitForKBQuestion(update tgbotapi.Update) (models.UserState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var answer string
	var err error

	doneChan := make(chan struct{}, 1)
	go func() {
		answer, err = controller.BackendClient.AskKB(ctx, update.Message.Text)
		doneChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		err = controller.Bot.SendMessage(
			update.Message.Chat.ID,
			"К сожалению не удалось получить ответ на Ваш вопрос",
		)
		if err != nil {
			return models.EmptyState, err
		}
	case <-doneChan:
		err = controller.Bot.SendMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("Ответ модели на Ваш вопрос: %s", answer),
		)
		if err != nil {
			return models.EmptyState, err
		}
	}

	return controller.EmptyState(update)
}
