package controllers

import (
	"context"
	"fmt"
	"time"
)

func (controller *Controller) SendStartupMessages(ctx context.Context) error {
	userIDs, err := controller.Storage.GetUserStatesUpdatedAfter(ctx, time.Now().Add(time.Minute*-40))
	if err != nil {
		return fmt.Errorf("failed to get user states: %w", err)
	}

	for _, userID := range userIDs {
		err := controller.Bot.SendMessage(userID, "Бот перезапущен")
		if err != nil {
			return fmt.Errorf("failed to send message to user: %w", err)
		}
	}

	return nil
}

func (controller *Controller) SendRestartMessages(ctx context.Context) error {
	userIDs, err := controller.Storage.GetUserStatesUpdatedAfter(ctx, time.Now().Add(time.Minute*-30))
	if err != nil {
		return fmt.Errorf("failed to get user states: %w", err)
	}

	for _, userID := range userIDs {
		err := controller.Bot.SendMessage(userID, "Бот перезапускается (ETA 2min)")
		if err != nil {
			return fmt.Errorf("failed to send message to user: %w", err)
		}
	}

	return nil
}
