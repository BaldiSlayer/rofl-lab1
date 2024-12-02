package mclient

import (
	"context"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

type ModelClient interface {
	// Ask отправляет запрос к модели
	Ask(ctx context.Context, question string, model string) (string, error)
	// AskWithContext отправляет запрос к модели с использованием контекста
	AskWithContext(ctx context.Context, question string, model string, questionContext []models.QAPair) (ResponseWithContext, error)
	// GetFormattedContext -
	GetFormattedContext(ctx context.Context, question string) ([]models.QAPair, error)
}

type ResponseWithContext struct {
	Answer  string
	Context string
}
