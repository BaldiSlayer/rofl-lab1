package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/githubclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/kbdatastorage"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
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
		answer, context, err := modelClient.AskWithContext(ctx, question, model)
		if err != nil {
			return KBAnswer{}, err
		}

		return KBAnswer{
			Model:   model,
			Answer:  answer,
			Context: &context,
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

func LoadQABase() ([]models.QAPair, error) {
	jsonStorage, err := kbdatastorage.NewJsonKBDataStorage("data/data.json")
	if err != nil {
		return nil, err
	}

	yamlStorage, err := kbdatastorage.NewYamlKBDataStorage("data/data.yaml")
	if err != nil {
		return nil, err
	}

	jsonData, err := jsonStorage.GetQAPairs()
	if err != nil {
		return nil, err
	}

	yamlData, err := yamlStorage.GetQAPairs()
	if err != nil {
		return nil, err
	}

	jsonData = []models.QAPair{} // NOTE: omit context from json temporarily

	slog.Info("Loaded context", "json", len(jsonData), "yaml", len(yamlData))

	return append(jsonData, yamlData...), nil
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
			Name:    fmt.Sprintf("%d-%s", i+2, answer.Model),
			Content: content,
		})
	}

	gist := githubclient.Gist{
		Description: version.BuildVersion(),
		Files:       files,
	}

	return ghClient.GistCreate(ctx, gist)
}
