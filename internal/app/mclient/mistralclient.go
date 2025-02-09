package mclient

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"log/slog"
	"net/http"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient/mistral"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

var _ ModelClient = (*Mistral)(nil)

type Mistral struct {
	*mistral.ClientWithResponses
}

const (
	llmServer = "http://llm-balancer"
	retryMax  = 5
)

func NewMistralClient() (ModelClient, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = retryMax
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

func (mc *Mistral) Ask(ctx context.Context, question string, model string) (string, error) {
	return mc.ask(ctx, question, nil, model)
}

func (mc *Mistral) AskWithContext(
	ctx context.Context,
	question string,
	model string,
	formattedContext string,
) (ResponseWithContext, error) {
	slog.Info("executing model request", "question", question, "context", formattedContext)

	answer, err := mc.ask(ctx, question, &formattedContext, model)

	return ResponseWithContext{
		Answer:      answer,
		Context:     formattedContext,
		ExtraPrompt: formattedContext,
	}, err
}

func (mc *Mistral) ask(ctx context.Context, question string, contextStr *string, model string) (string, error) {
	resp, err := mc.ApiGetChatResponseGetChatResponsePostWithResponse(ctx, mistral.GetChatResponseRequest{
		Context: contextStr,
		Model:   &model,
		Prompt:  question,
	})

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("error while requesting LLM: status code: %d", resp.StatusCode())
	}

	return resp.JSON200.Response, nil
}

func (mc *Mistral) processQuestionsRequest(ctx context.Context, question string) ([]mistral.QuestionAnswer, error) {
	resp, err := mc.ApiProcessQuestionsProcessQuestionsPostWithResponse(ctx, mistral.SearchSimilarRequest{
		Question: question,
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error requesting LLM: status code %d", resp.StatusCode())
	}

	return resp.JSON200.Result, nil
}

func (mc *Mistral) GetFormattedContext(ctx context.Context, question string) ([]models.QAPair, error) {
	contextQASlice, err := mc.processQuestionsRequest(ctx, question)
	if err != nil {
		return nil, err
	}

	qaPairSlice := make([]models.QAPair, 0, len(contextQASlice))
	for _, val := range contextQASlice {
		qaPairSlice = append(qaPairSlice, models.QAPair{
			Question: val.Question,
			Answer:   val.Answer,
		})
	}

	return qaPairSlice, nil
}
