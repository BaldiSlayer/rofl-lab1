package ustorage

import (
	"errors"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

// MapUserStorage хранилище данных о пользователе, которое использует map
type MapUserStorage struct {
	states map[int64]models.UserState
	formalTRS map[int64]string
	TRS map[int64]trsparser.Trs
	reqeusts map[int64]string
	parseErrors map[int64]string
}

func NewMapUserStorage() (*MapUserStorage, error) {
	return &MapUserStorage{
		states: make(map[int64]models.UserState),
		formalTRS: make(map[int64]string),
		TRS: make(map[int64]trsparser.Trs),
		reqeusts: make(map[int64]string),
		parseErrors: make(map[int64]string),
	}, nil
}

func (s *MapUserStorage) GetState(userID int64) (models.UserState, error) {
	return s.states[userID], nil
}

func (s *MapUserStorage) GetTRS(userID int64) (trsparser.Trs, error) {
	if trs, ok := s.TRS[userID]; ok {
		return trs, nil
	}
	return trsparser.Trs{}, errors.New("trs not found")
}

func (s *MapUserStorage) GetFormalTRS(userID int64) (string, error) {
	if trs, ok := s.formalTRS[userID]; ok {
		return trs, nil
	}
	return "", errors.New("formal trs not found")
}

func (s *MapUserStorage) GetRequest(userID int64) (string, error) {
	if request, ok := s.reqeusts[userID]; ok {
		return request, nil
	}
	return "", errors.New("user request not found")
}

func (s *MapUserStorage) GetParseError(userID int64) (string, error) {
	if parseError, ok := s.parseErrors[userID]; ok {
		return parseError, nil
	}
	return "", errors.New("parse error not found")
}

func (s *MapUserStorage) SetState(userID int64, state models.UserState) error {
	s.states[userID] = state
	return nil
}

func (s *MapUserStorage) SetTRS(userID int64, trs trsparser.Trs) error {
	s.TRS[userID] = trs
	return nil
}

func (s *MapUserStorage) SetFormalTRS(userID int64, formalTrs string) error {
	s.formalTRS[userID] = formalTrs
	return nil
}

func (s *MapUserStorage) SetRequest(userID int64, request string) error {
	s.reqeusts[userID] = request
	return nil
}

func (s *MapUserStorage) SetParseError(userID int64, parseError string) error {
	s.parseErrors[userID] = parseError
	return nil
}
