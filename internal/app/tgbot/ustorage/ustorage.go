package ustorage

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

type UserDataStorage interface {
	SetState(userID int64, state models.UserState) error // +
	SetTRS(userID int64, trs trsparser.Trs) error
	SetFormalTRS(userID int64, formalTrs string) error
	SetRequest(userID int64, request string) error       // +
	SetParseError(userID int64, parseError string) error // +
	GetState(userID int64) (models.UserState, error)
	GetTRS(userID int64) (trsparser.Trs, error)
	GetFormalTRS(userID int64) (string, error)
	GetRequest(userID int64) (string, error)
	GetParseError(userID int64) (string, error)
}
