package mclient

import (
	"context"
)

type ModelClient interface {
	// Ask отправляет запрос к модели
	Ask(ctx context.Context, question string, model string) (string, error)
	// AskWithContext отправляет запрос к модели с использованием контекста
	AskWithContext(ctx context.Context, question string, model string) (ResponseWithContext, error)
}

type ResponseWithContext struct {
	Answer  string
	Context string
}
