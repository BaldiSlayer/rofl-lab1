package usecases

import (
	"context"
	"log/slog"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/githubclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/kbdatastorage"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/version"
)

func AskKnowledgeBase(ctx context.Context, modelClient mclient.ModelClient, question string) (string, string, error) {
	return modelClient.AskWithContext(ctx, question, "open-mistral-7b")
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

	slog.Info("Loaded context", "json", len(jsonData), "yaml", len(yamlData))

	return append(jsonData, yamlData...), nil
}

func UploadKnowledgeBaseAnswer(ctx context.Context, ghClient *githubclient.Client, userRequest, usedContext, answer string) (string, error) {
	gist := githubclient.Gist{
		Description: version.BuildVersion(),
		Files: []githubclient.GistFile{
			{
				Name:    "1-question.md",
				Content: userRequest,
			},
			{
				Name:    "2-context.md",
				Content: usedContext,
			},
			{
				Name:    "3-answer.md",
				Content: answer,
			},
		},
	}

	return ghClient.GistCreate(ctx, gist)
}
