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

func getContextFromQASlice(contextQASlice []models.QAPair) string {
	result := ""

	for _, item := range contextQASlice {
		result = result + fmt.Sprintf("\nВопрос: %s Ответ: %s", item.Question, item.Answer)
	}

	return result
}

func (mc *Mistral) AskWithContext(
	ctx context.Context,
	question string,
	model string,
	questionContext []models.QAPair,
) (ResponseWithContext, error) {
	formattedContext := getContextFromQASlice(questionContext)

	slog.Info("executing model request", "question", question, "context", formattedContext)

	answer, err := mc.ask(ctx, question, &formattedContext, model)
	return ResponseWithContext{
		Answer:  answer,
		Context: formattedContext,
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
