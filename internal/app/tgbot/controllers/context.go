package controllers

import (
	"context"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/usecases"
)

func (controller *Controller) InitContext(ctx context.Context) error {
	qa, err := usecases.LoadQABase()
	if err != nil {
		return err
	}
	return controller.ModelClient.InitContext(ctx, qa)
}
