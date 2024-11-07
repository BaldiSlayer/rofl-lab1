package usecases

import (
	"context"
	"log/slog"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/kbdatastorage"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

func AskKnowledgeBase(ctx context.Context, modelClient mclient.ModelClient, question string) (string, error) {
	return modelClient.AskWithContext(ctx, question)
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
