package ustorage

import "github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"

type UserDataStorage interface {
	SetState(userID int64, state models.UserState) error
	GetState(userID int64) models.UserState
}
