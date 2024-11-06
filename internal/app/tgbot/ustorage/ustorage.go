package ustorage

import (
	"context"
	"errors"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

type UserDataStorage interface {
	SetState(ctx context.Context, userID int64, state models.UserState) error
	SetTRS(ctx context.Context, userID int64, trs trsparser.Trs) error
	SetFormalTRS(ctx context.Context, userID int64, formalTrs string) error
	SetRequest(ctx context.Context, userID int64, request string) error
	SetParseError(ctx context.Context, userID int64, parseError string) error
	GetState(ctx context.Context, userID int64) (models.UserState, error)
	GetTRS(ctx context.Context, userID int64) (trsparser.Trs, error)
	GetFormalTRS(ctx context.Context, userID int64) (string, error)
	GetRequest(ctx context.Context, userID int64) (string, error)
	GetParseError(ctx context.Context, userID int64) (string, error)
}

var ErrNotFound = errors.New("not found")
