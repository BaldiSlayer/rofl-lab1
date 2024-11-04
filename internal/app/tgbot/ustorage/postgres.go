package ustorage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

// PostgresUserStorage хранилище данных о пользователе, которое использует map
type PostgresUserStorage struct {
	pg *pgxpool.Pool
}

func NewPostgresUserStorage() (*PostgresUserStorage, error) {
	connString := fmt.Sprintf("user=%s password=%s host=postgres "+
		"port=5432 dbname=%s pool_max_conns=10",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &PostgresUserStorage{
		pg: dbpool,
	}, nil
}

func (s *PostgresUserStorage) GetState(userID int64) (models.UserState, error) {
	var state int
	err := s.pg.QueryRow(context.TODO(), "SELECT state FROM tfllab1.user_state WHERE user_id = $1", userID).Scan(&state)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}

	return models.UserState(state), nil
}

func (s *PostgresUserStorage) GetTRS(userID int64) (trsparser.Trs, error) {
	var trs trsparser.Trs
	err := s.pg.QueryRow(context.TODO(), "SELECT parse_result FROM tfllab1.extraction_result WHERE user_id = $1", userID).Scan(&trs)
	if errors.Is(err, pgx.ErrNoRows) {
		return trsparser.Trs{}, ErrNotFound
	}
	if err != nil {
		return trsparser.Trs{}, err
	}

	return trs, nil
}

func (s *PostgresUserStorage) GetFormalTRS(userID int64) (string, error) {
	var trs string
	err := s.pg.QueryRow(context.TODO(), "SELECT formalize_result FROM tfllab1.extraction_result WHERE user_id = $1", userID).Scan(&trs)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return trs, nil
}

func (s *PostgresUserStorage) GetRequest(userID int64) (string, error) {
	var userRequest string
	err := s.pg.QueryRow(context.TODO(), "SELECT user_request FROM tfllab1.extraction_result WHERE user_id = $1", userID).Scan(&userRequest)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return userRequest, nil
}

func (s *PostgresUserStorage) GetParseError(userID int64) (string, error) {
	var parseError string
	err := s.pg.QueryRow(context.TODO(), "SELECT parse_error FROM tfllab1.extraction_result WHERE user_id = $1", userID).Scan(&parseError)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return parseError, nil
}

func (s *PostgresUserStorage) SetState(userID int64, state models.UserState) error {
	_, err := s.pg.Exec(context.TODO(), "INSERT INTO tfllab1.user_state(user_id, state) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET state=excluded.state", userID, state)
	return err
}

func (s *PostgresUserStorage) SetTRS(userID int64, trs trsparser.Trs) error {
	_, err := s.pg.Exec(context.TODO(), "INSERT INTO tfllab1.extraction_result(user_id, parse_result) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET parse_result=EXCLUDED.parse_result", userID, trs)
	return err
}

func (s *PostgresUserStorage) SetFormalTRS(userID int64, formalTrs string) error {
	_, err := s.pg.Exec(context.TODO(), "INSERT INTO tfllab1.extraction_result(user_id, formalize_result) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET formalize_result=EXCLUDED.formalize_result", userID, formalTrs)
	return err
}

func (s *PostgresUserStorage) SetRequest(userID int64, request string) error {
	_, err := s.pg.Exec(context.TODO(), "INSERT INTO tfllab1.extraction_result(user_id, user_request) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET user_request=EXCLUDED.user_request", userID, request)
	return err
}

func (s *PostgresUserStorage) SetParseError(userID int64, parseError string) error {
	_, err := s.pg.Exec(context.TODO(), "INSERT INTO tfllab1.extraction_result(user_id, parse_error) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET parse_error=EXCLUDED.parse_error", userID, parseError)
	return err
}
