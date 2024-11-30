package usecases

import (
	"context"
	"fmt"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/githubclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/version"
)

type KBAnswer struct {
	Model   string
	Answer  string
	Context *string
}

func AskKnowledgeBase(ctx context.Context, modelClient mclient.ModelClient, question string) ([]KBAnswer, error) {
	res := make([]KBAnswer, 0)

	requests := []struct {
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

	for _, request := range requests {
		var formattedContext string
		var err error

		if request.useContext {
			formattedContext, err = modelClient.GetFormattedContext(ctx, question)
			if err != nil {
				return nil, err
			}
		}

		ans, err := ask(ctx, modelClient, question, request.model, formattedContext)
		if err != nil {
			return nil, err
		}

		res = append(res, ans)
	}

	return res, nil
}

func ask(
	ctx context.Context,
	modelClient mclient.ModelClient,
	question,
	model string,
	formattedContext string,
) (KBAnswer, error) {
	if formattedContext != "" {
		res, err := modelClient.AskWithContext(ctx, question, model, formattedContext)
		if err != nil {
			return KBAnswer{}, fmt.Errorf("failed to ask model with context: %w", err)
		}

		return KBAnswer{
			Model:   model,
			Answer:  res.Answer,
			Context: &res.Context,
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

func UploadKnowledgeBaseAnswers(
	ctx context.Context,
	ghClient *githubclient.Client,
	userRequest string,
	answers []KBAnswer,
) (string, error) {
	files := make([]githubclient.GistFile, 0)

	// предполагаем, что контекст у всех ответов одинаковый
	for _, answer := range answers {
		if answer.Context != nil {
			files = append(files, githubclient.GistFile{
				Name:    "context.md",
				Content: *answer.Context,
			})

			break
		}
	}

	files = append(files, githubclient.GistFile{
		Name:    "1-question.md",
		Content: userRequest,
	})

	for i, answer := range answers {
		extraInfo := fmt.Sprintf("## %s. Контекст был использован", answer.Model)

		if answer.Context == nil {
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
