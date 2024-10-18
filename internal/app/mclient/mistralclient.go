package mclient

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient/mistral"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

var _ ModelClient = (*Mistral)(nil)

type Mistral struct {
	*mistral.ClientWithResponses
}

// TODO: configure?
const LlmServer = "http://llm:8100"

func NewMistralClient(questions []models.QAPair) (ModelClient, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	standardClient := retryClient.StandardClient() // *http.Client

	c, err := mistral.NewClientWithResponses(LlmServer, mistral.WithHTTPClient(standardClient))
	if err != nil {
		return nil, err
	}

	mc := &Mistral{
		ClientWithResponses: c,
	}

	message, err := mc.processQuestionsRequest(questions, false)
	if err != nil {
		return nil, err
	}

	slog.Info("Initialized llm context", "messages", message)

	return mc, nil
}

func (mc *Mistral) Ask(question string) (string, error) {
	return mc.ask(question, nil)
}

func (mc *Mistral) AskWithContext(question string) (string, error) {
	context, err := mc.processQuestionsRequest([]models.QAPair{{
		Question: question,
		Answer:   "",
	}}, true)
	if err != nil {
		return "", err
	}

	if len(context) != 1 {
		return "", fmt.Errorf("expected single answer from process_questions endpoint, got %d", len(context))
	}

	return mc.ask(question, &context[0])
}

func (mc *Mistral) ask(question string, contextStr *string) (string, error) {
	resp, err := mc.ApiGetChatResponseGetChatResponsePostWithResponse(context.TODO(), mistral.GetChatResponseRequest{
		Context: contextStr,
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

func toQuestionsList(QAPairs []models.QAPair) []mistral.QuestionAnswer {
	res := []mistral.QuestionAnswer{}
	for _, qa := range QAPairs {
		res = append(res, mistral.QuestionAnswer{
			Answer:   qa.Answer,
			Question: qa.Question,
		})
	}
	return res
}

func (mc *Mistral) processQuestionsRequest(QAPairs []models.QAPair, useSaved bool) ([]string, error) {
	resp, err := mc.ApiProcessQuestionsProcessQuestionsPostWithResponse(context.TODO(), mistral.ProcessQuestionsRequest{
		Filename:      nil,
		QuestionsList: toQuestionsList(QAPairs),
		UseSaved:      &useSaved,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		slog.Error("error requesting LLM", "code", resp.StatusCode())
		return nil, errors.New("error requesting LLM")
	}

	return resp.JSON200.Result, nil
}
