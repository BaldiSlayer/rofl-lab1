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
	res := []KBAnswer{}
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
		ans, err := ask(ctx, modelClient, question, request.model, request.useContext)
		if err != nil {
			return nil, err
		}
		res = append(res, ans)
	}

	return res, nil
}

func ask(ctx context.Context, modelClient mclient.ModelClient, question, model string, useContext bool) (KBAnswer, error) {
	if useContext {
		res, err := modelClient.AskWithContext(ctx, question, model)
		if err != nil {
			return KBAnswer{}, err
		}

		return KBAnswer{
			Model:   model,
			Answer:  res.Answer,
			Context: &res.Context,
		}, nil
	}

	answer, err := modelClient.Ask(ctx, question, model)
	if err != nil {
		return KBAnswer{}, err
	}

	return KBAnswer{
		Model:   model,
		Answer:  answer,
		Context: nil,
	}, nil
}

func UploadKnowledgeBaseAnswers(ctx context.Context, ghClient *githubclient.Client, userRequest string, answers []KBAnswer) (string, error) {
	files := []githubclient.GistFile{}
	files = append(files, githubclient.GistFile{
		Name:    "1-question.md",
		Content: userRequest,
	})
	for i, answer := range answers {
		var content string
		if answer.Context == nil {
			content = fmt.Sprintf("%s", answer.Answer)
		} else {
			content = fmt.Sprintf("%s\n\nИспользован контекст:\n%s", answer.Answer, *answer.Context)
		}

		files = append(files, githubclient.GistFile{
			Name:    fmt.Sprintf("%d-%s.md", i+2, answer.Model),
			Content: content,
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
