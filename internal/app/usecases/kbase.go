package usecases

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
)

func AskKnowledgeBase(modelClient mclient.ModelClient, question string) (string, error) {
	return modelClient.AskWithContext(question)
}
