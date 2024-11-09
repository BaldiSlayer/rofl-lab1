package kbdatastorage

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/models"
)

type KBDataStorage interface {
	GetQAPairs([]models.QAPair, error)
}
