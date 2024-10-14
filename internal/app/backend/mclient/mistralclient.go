package mclient

import "github.com/BaldiSlayer/rofl-lab1/internal/app/backend/models"

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../../../docs/llm-api.yaml

var _ ModelClient = (*Mistral)(nil)

type Mistral struct{}

func (mc *Mistral) Ask(question string) (string, error) {
	return question, nil
}

func (mc *Mistral) AskWithContext(question string, answerContext []models.QAPair) (string, error) {
	return question, nil
}
