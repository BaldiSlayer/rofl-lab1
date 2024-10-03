package vdatabase

import "github.com/BaldiSlayer/rofl-lab1/internal/app/backend/models"

type VectorDatabase interface {
	GetSimilar(question string) ([]models.QAPair, error)
}
