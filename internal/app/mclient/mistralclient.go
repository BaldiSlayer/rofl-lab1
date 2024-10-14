package mclient

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient/mistral"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

var _ ModelClient = (*Mistral)(nil)

type Mistral struct{
	*mistral.ClientWithResponses
}

const LlmServer = "http://llm:8100"

func NewMistralClient() (ModelClient, error) {
	c, err := mistral.NewClientWithResponses(LlmServer)
	if err != nil {
		return nil, err
	}

	return &Mistral{
		ClientWithResponses: c,
	}, nil
}

func (mc *Mistral) Ask(question string) (string, error) {
	resp, err := mc.ApiGetChatResponseGetChatResponsePostWithResponse(context.TODO(), mistral.GetChatResponseRequest{
		Context: nil,
		Model:   nil,
		Prompt:  question,
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != http.StatusOK {
		slog.Error("error requesting LLM", "code", resp.StatusCode())
		return "", errors.New("error requesting LLM")
	}

	return resp.JSON200.Response, nil
}

func (mc *Mistral) AskWithContext(question string, answerContext []models.QAPair) (string, error) {
	return question, nil
}
