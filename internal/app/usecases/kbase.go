package usecases

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"strings"
	"sync"

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
	FormattedContext string
}

func getContextFromQASlice(contextQASlice []models.QAPair) (string, error) {
	t, err := template.New("qaTemplate").Parse(ModelContextTemplate)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	err = t.Execute(&output, contextQASlice)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func GetSimilarElements(
	ctx context.Context,
	modelClient mclient.ModelClient,
	ghClient *githubclient.Client,
	question string,
) (string, error) {
	questionsContext, err := modelClient.GetFormattedContext(ctx, question)
	if err != nil {
		return "", fmt.Errorf("failed to get formatted context: %w", err)
	}

	files := []githubclient.GistFile{
		{
			Name:    "1-question.md",
			Content: fmt.Sprintf("## Пользовательский вопрос\n%s", question),
		},
		{
			Name:    fmt.Sprintf("similar.md"),
			Content: getContextPresentationForGist(questionsContext, "", false),
		},
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

func AskKnowledgeBase(
	ctx context.Context,
	modelClient mclient.ModelClient,
	question string,
	requests []commons.ModelRequest,
) (AskResults, error) {
	res := make([]KBAnswer, len(requests))

	questionsContext, err := modelClient.GetFormattedContext(ctx, question)
	if err != nil {
		return AskResults{}, fmt.Errorf("failed to get formatted context: %w", err)
	}

	formattedContext, err := getContextFromQASlice(questionsContext)
	if err != nil {
		return AskResults{}, fmt.Errorf("failed to gen template for extra prompt: %w", err)
	}

	var wg sync.WaitGroup

	wg.Add(len(requests))

	for i, request := range requests {
		go func(i int, request models.ModelRequest) {
			defer wg.Done()

			questionContext := ""

			if request.UseContext {
				questionContext = formattedContext
			}

			ans, err := ask(ctx, modelClient, question, request.Model, request.UseContext, questionContext)
			if err != nil {
				slog.Error("error while asking model", "model", request.Model)

				return
			}

			res[i] = ans
		}(i, request)
	}

	wg.Wait()

	return AskResults{
		Answers:          res,
		QuestionsContext: questionsContext,
		FormattedContext: formattedContext,
	}, nil
}

func ask(
	ctx context.Context,
	modelClient mclient.ModelClient,
	question string,
	model string,
	useContext bool,
	formattedContext string,
) (KBAnswer, error) {
	if useContext {
		res, err := modelClient.AskWithContext(
			ctx,
			question,
			model,
			formattedContext,
		)
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

func getContextPresentationForGist(
	questionsContext []models.QAPair,
	extraPrompt string,
	withPrompt bool,
) string {
	var sb strings.Builder

	if withPrompt {
		sb.WriteString("## Отправленный промпт с контекстом \n")
		sb.WriteString("```")
		sb.WriteString(extraPrompt)
		sb.WriteString("\n```\n")

		sb.WriteString("## Контекст")
		sb.WriteByte('\n')
	}

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
			Name: "context.md",
			Content: getContextPresentationForGist(
				askResults.QuestionsContext,
				askResults.FormattedContext,
				true,
			),
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
