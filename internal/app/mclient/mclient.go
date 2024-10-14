package mclient

import "github.com/BaldiSlayer/rofl-lab1/internal/app/models"

type ModelClient interface {
	// Ask отправляет запрос к модели
	Ask(question string) (string, error)
	// AskWithContext отправляет запрос к модели с дополнительным контекстом
	AskWithContext(question string, answerContext []models.QAPair) (string, error)
}
