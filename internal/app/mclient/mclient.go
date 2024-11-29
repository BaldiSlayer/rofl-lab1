package mclient

import (
	"context"
)

type ModelClient interface {
	// Ask отправляет запрос к модели
	Ask(ctx context.Context, question string, model string) (string, error)
	// AskWithContext отправляет запрос к модели с использованием контекста
	AskWithContext(ctx context.Context, question string, model string, formattedContext string) (ResponseWithContext, error)
	// GetFormattedContext -
	GetFormattedContext(ctx context.Context, question string) (string, error)
}

type ResponseWithContext struct {
	Answer  string
	Context string
}
