package ustorage

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"sync"
)

// MapUserStorage хранилище данных о пользователе, которое использует map
type MapUserStorage struct {
	storage map[int64]models.UserState
	mu      sync.Mutex
}

func NewMapUserStorage() (*MapUserStorage, error) {
	return &MapUserStorage{
		storage: make(map[int64]models.UserState),
	}, nil
}

func (s *MapUserStorage) SetState(userID int64, state models.UserState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[userID] = state

	return nil
}

func (s *MapUserStorage) GetState(userID int64) models.UserState {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.storage[userID]
}
