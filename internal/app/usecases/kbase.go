package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/githubclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
	commons "github.com/BaldiSlayer/rofl-lab1/internal/app/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/version"
)

type KBAnswer struct {
	Model      string
	Answer     string
	UseContext bool
}

type AskResults struct {
	Answers          []KBAnswer
	QuestionsContext []models.QAPair
}

func AskKnowledgeBase(
	ctx context.Context,
	modelClient mclient.ModelClient,
	question string,
	requests []commons.ModelRequest,
) (AskResults, error) {
	res := make([]KBAnswer, 0, len(requests))

	questionsContext, err := modelClient.GetFormattedContext(ctx, question)
	if err != nil {
		return AskResults{}, fmt.Errorf("failed to get formatted context: %w", err)
	}

	for _, request := range requests {
		questionContext := []models.QAPair(nil)

		if request.UseContext {
			questionContext = questionsContext
		}

		ans, err := ask(ctx, modelClient, question, request.Model, questionContext)
		if err != nil {
			return AskResults{}, fmt.Errorf("error while ask model: %w", err)
		}

		res = append(res, ans)
	}

	return AskResults{
		Answers:          res,
		QuestionsContext: questionsContext,
	}, nil
}

func ask(
	ctx context.Context,
	modelClient mclient.ModelClient,
	question,
	model string,
	questionContext []models.QAPair,
) (KBAnswer, error) {
	if len(questionContext) != 0 {
		res, err := modelClient.AskWithContext(ctx, question, model, questionContext)
		if err != nil {
			return KBAnswer{}, fmt.Errorf("failed to ask model with context: %w", err)
		}

		return KBAnswer{
			Model:      model,
			Answer:     res.Answer,
			UseContext: true,
		}, nil
	}

	answer, err := modelClient.Ask(ctx, question, model)
	if err != nil {
		return KBAnswer{}, fmt.Errorf("failed to ask model without context: %w", err)
	}

	return KBAnswer{
		Model:  model,
		Answer: answer,
	}, nil
}

func getContextPresentationForGist(questionsContext []models.QAPair) string {
	var sb strings.Builder

	sb.WriteString("## Шаблон промпта для контекста\n")
	sb.WriteString("```")
	sb.WriteString(mclient.ModelContextTemplate)
	sb.WriteString("\n```\n")

	sb.WriteString("## Контекст")
	sb.WriteByte('\n')

	for _, val := range questionsContext {
		sb.WriteString("### Вопрос\n")
		sb.WriteString(val.Question)
		sb.WriteByte('\n')
		sb.WriteString("### Ответ\n")
		sb.WriteString(val.Answer)
		sb.WriteByte('\n')
	}

	return sb.String()
}

func UploadKnowledgeBaseAnswers(
	ctx context.Context,
	ghClient *githubclient.Client,
	userRequest string,
	askResults AskResults,
) (string, error) {
	files := make([]githubclient.GistFile, 0, 4)

	files = append(
		files,
		githubclient.GistFile{
			Name:    "context.md",
			Content: getContextPresentationForGist(askResults.QuestionsContext),
		},
		githubclient.GistFile{
			Name:    "1-question.md",
			Content: fmt.Sprintf("## Пользовательский вопрос\n%s", userRequest),
		},
	)

	for i, answer := range askResults.Answers {
		extraInfo := fmt.Sprintf("## %s. Контекст был использован", answer.Model)

		if !answer.UseContext {
			extraInfo = fmt.Sprintf("## %s. Контекст не был использован", answer.Model)
		}

		files = append(files, githubclient.GistFile{
			Name:    fmt.Sprintf("%d-%s.md", i+2, answer.Model),
			Content: fmt.Sprintf("%s\n%s", extraInfo, answer.Answer),
		})
	}

	gist := githubclient.Gist{
		Description: version.BuildVersion(),
		Files:       files,
	}

	link, err := ghClient.GistCreate(ctx, gist)
	if err != nil {
		return "", fmt.Errorf("failed to create gist: %w", err)
	}

	return link, nil
}
