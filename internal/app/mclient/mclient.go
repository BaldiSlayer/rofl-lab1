package mclient

import (
	"context"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

type ModelClient interface {
	// InitContext отправляет запрос инициализации контекста
	InitContext(ctx context.Context, data []models.QAPair) error
	// Ask отправляет запрос к модели
	Ask(ctx context.Context, question string) (string, error)
	// AskWithContext отправляет запрос к модели с использованием контекста
	AskWithContext(ctx context.Context, question string) (answer string, context string, err error)
}
