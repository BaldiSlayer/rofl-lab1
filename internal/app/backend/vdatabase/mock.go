package vdatabase

import "github.com/BaldiSlayer/rofl-lab1/internal/app/backend/models"

// CMem - реализация векторной базы данных с помощью chromem-go
type Mock struct {
}

func (cm *Mock) Init(baseURL string, contextFile string) error {
	return nil
}

func (cm *Mock) GetSimilar(question string) ([]models.QAPair, error) {
	return []models.QAPair{
		{
			Question: question,
			Answer:   "idk",
		},
	}, nil
}
