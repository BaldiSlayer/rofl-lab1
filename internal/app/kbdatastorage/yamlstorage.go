package kbdatastorage

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

type YamlKBDataStorage struct {
	path string
}

func NewYamlKBDataStorage(path string) (*YamlKBDataStorage, error) {
	return &YamlKBDataStorage{
		path: path,
	}, nil
}

func (s *YamlKBDataStorage) GetQAPairs() ([]models.QAPair, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var res []models.QAPair

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&res)

	return res, err
}
