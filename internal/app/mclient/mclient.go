package mclient

import (
	"context"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
	commons "github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

func GetDefaultModelRequestsPattern() []commons.ModelRequest {
	return []commons.ModelRequest{
		{
			Model:      "mistral-large-2411",
			UseContext: true,
		},
		{
			Model:      "mistral-large-2411",
			UseContext: false,
		},
		{
			Model:      "open-mistral-7b",
			UseContext: true,
		},
	}
}

func GetFastModelRequestsPattern() []commons.ModelRequest {
	return []commons.ModelRequest{
		{
			Model:      "mistral-large-2411",
			UseContext: true,
		},
	}
}

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
