package usecases

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/kbdatastorage"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

func AskKnowledgeBase(modelClient mclient.ModelClient, question string) (string, error) {
	return modelClient.AskWithContext(question)
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

	return append(jsonData, yamlData...), nil
}
