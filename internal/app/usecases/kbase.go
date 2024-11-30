package usecases

import (
	"context"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
	"strings"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/githubclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
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

func AskKnowledgeBase(ctx context.Context, modelClient mclient.ModelClient, question string) (AskResults, error) {
	res := make([]KBAnswer, 0)

	requests := [...]struct {
		model      string
		useContext bool
	}{
		{
			model:      "mistral-large-latest",
			useContext: true,
		},
		{
			model:      "mistral-large-latest",
			useContext: false,
		},
		{
			model:      "open-mistral-7b",
			useContext: true,
		},
	}

	questionsContext, err := modelClient.GetFormattedContext(ctx, question)
	if err != nil {
		return AskResults{}, err
	}

	for _, request := range requests {
		questionContext := questionsContext

		// если контекст не используем, то делаем слайс ниловым
		if !request.useContext {
			questionContext = []models.QAPair(nil)
		}

		ans, err := ask(ctx, modelClient, question, request.model, questionContext)
		if err != nil {
			return AskResults{}, err
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

	files = append(files, githubclient.GistFile{
		Name:    "context.md",
		Content: getContextPresentationForGist(askResults.QuestionsContext),
	})

	files = append(files, githubclient.GistFile{
		Name:    "1-question.md",
		Content: userRequest,
	})

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
