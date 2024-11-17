package kbdatastorage

import (
	"encoding/json"
	"os"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

type JsonKBDataStorage struct {
	path string
}

func NewJsonKBDataStorage(path string) (*JsonKBDataStorage, error) {
	return &JsonKBDataStorage{
		path: path,
	}, nil
}

func (s *JsonKBDataStorage) GetQAPairs() ([]models.QAPair, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var res []models.QAPair

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&res)

	return res, err
}
