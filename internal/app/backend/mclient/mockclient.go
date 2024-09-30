package mclient

import "github.com/BaldiSlayer/rofl-lab1/internal/app/backend/models"

type Mock struct{}

func (mc *Mock) Ask(question string) (string, error) {
	return question, nil
}

func (mc *Mock) AskWithContext(question string, answerContext []models.QAPair) (string, error) {
	return question, nil
}
