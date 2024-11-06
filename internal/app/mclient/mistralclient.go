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

const llmServer = "http://llm:8100"

func NewMistralClient(questions []models.QAPair) (ModelClient, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	standardClient := retryClient.StandardClient() // *http.Client

	c, err := mistral.NewClientWithResponses(llmServer, mistral.WithHTTPClient(standardClient))
	if err != nil {
		return nil, err
	}

	mc := &Mistral{
		ClientWithResponses: c,
	}

	return mc, nil
}

func (mc *Mistral) InitContext(ctx context.Context, questions []models.QAPair) error {
	_, err := mc.processQuestionsRequest(ctx, questions, false)
	return err
}

func (mc *Mistral) Ask(ctx context.Context, question string) (string, error) {
	return mc.ask(ctx, question, nil)
}

func (mc *Mistral) AskWithContext(ctx context.Context, question string) (string, error) {
	contexts, err := mc.processQuestionsRequest(ctx, []models.QAPair{{
		Question: question,
		Answer:   "",
	}}, true)
	if err != nil {
		return "", err
	}

	if len(contexts) != 1 {
		return "", fmt.Errorf("expected single answer from process_questions endpoint, got %d", len(contexts))
	}

	context := contexts[0]

	slog.Info("executing model request", "question", question, "context", context)

	return mc.ask(ctx, question, &context)
}

func (mc *Mistral) ask(ctx context.Context, question string, contextStr *string) (string, error) {
	resp, err := mc.ApiGetChatResponseGetChatResponsePostWithResponse(ctx, mistral.GetChatResponseRequest{
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

func (mc *Mistral) processQuestionsRequest(ctx context.Context, QAPairs []models.QAPair, useSaved bool) ([]string, error) {
	resp, err := mc.ApiProcessQuestionsProcessQuestionsPostWithResponse(ctx, mistral.ProcessQuestionsRequest{
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
