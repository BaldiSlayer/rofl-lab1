package ustorage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"
)

// PostgresStorage хранилище данных о пользователе, которое использует Postgres
type PostgresStorage struct {
	pg *pgxpool.Pool
}

func NewPostgresStorage(user, password, db string) (*PostgresStorage, error) {
	connString := fmt.Sprintf("user=%s password=%s host=postgres "+
		"port=5432 dbname=%s pool_max_conns=10",
		user,
		password,
		db,
	)

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		pg: dbpool,
	}, nil
}

func (s *PostgresStorage) Close() {
	s.pg.Close()
}

func (s *PostgresStorage) GetState(ctx context.Context, userID int64) (models.UserState, error) {
	var state int
	err := s.pg.QueryRow(ctx, "SELECT state FROM tfllab1.user_state WHERE user_id = $1", userID).Scan(&state)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}

	return models.UserState(state), nil
}

func (s *PostgresStorage) GetTRS(ctx context.Context, userID int64) (trsparser.Trs, error) {
	var trs trsparser.Trs
	err := s.pg.QueryRow(ctx, "SELECT parse_result FROM tfllab1.extraction_result WHERE user_id = $1", userID).Scan(&trs)
	if errors.Is(err, pgx.ErrNoRows) {
		return trsparser.Trs{}, ErrNotFound
	}
	if err != nil {
		return trsparser.Trs{}, err
	}

	return trs, nil
}

func (s *PostgresStorage) GetFormalTRS(ctx context.Context, userID int64) (string, error) {
	var trs string
	err := s.pg.QueryRow(ctx, "SELECT formalize_result FROM tfllab1.extraction_result WHERE user_id = $1", userID).Scan(&trs)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return trs, nil
}

func (s *PostgresStorage) GetRequest(ctx context.Context, userID int64) (string, error) {
	var userRequest string
	err := s.pg.QueryRow(ctx, "SELECT user_request FROM tfllab1.extraction_result WHERE user_id = $1", userID).Scan(&userRequest)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return userRequest, nil
}

func (s *PostgresStorage) GetParseError(ctx context.Context, userID int64) (string, error) {
	var parseError string
	err := s.pg.QueryRow(ctx, "SELECT parse_error FROM tfllab1.extraction_result WHERE user_id = $1", userID).Scan(&parseError)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return parseError, nil
}

func (s *PostgresStorage) SetState(ctx context.Context, userID int64, state models.UserState) error {
	_, err := s.pg.Exec(ctx, "INSERT INTO tfllab1.user_state(user_id, state) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET state=excluded.state", userID, state)
	return err
}

func (s *PostgresStorage) SetTRS(ctx context.Context, userID int64, trs trsparser.Trs) error {
	_, err := s.pg.Exec(ctx, "INSERT INTO tfllab1.extraction_result(user_id, parse_result) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET parse_result=EXCLUDED.parse_result", userID, trs)
	return err
}

func (s *PostgresStorage) SetFormalTRS(ctx context.Context, userID int64, formalTrs string) error {
	_, err := s.pg.Exec(ctx, "INSERT INTO tfllab1.extraction_result(user_id, formalize_result) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET formalize_result=EXCLUDED.formalize_result", userID, formalTrs)
	return err
}

func (s *PostgresStorage) SetRequest(ctx context.Context, userID int64, request string) error {
	_, err := s.pg.Exec(ctx, "INSERT INTO tfllab1.extraction_result(user_id, user_request) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET user_request=EXCLUDED.user_request", userID, request)
	return err
}

func (s *PostgresStorage) SetParseError(ctx context.Context, userID int64, parseError string) error {
	_, err := s.pg.Exec(ctx, "INSERT INTO tfllab1.extraction_result(user_id, parse_error) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET parse_error=EXCLUDED.parse_error", userID, parseError)
	return err
}

func (s *PostgresStorage) GetUserStatesUpdatedAfter(ctx context.Context, after time.Time) ([]int64, error) {
	rows, err := s.pg.Query(ctx, "SELECT user_id FROM tfllab1.user_state WHERE updated_at > $1", after)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userIDs := []int64{}

	for rows.Next() {
		var userID int64
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (s *PostgresStorage) TryLock(ctx context.Context, userID int64, instanceID uuid.UUID, duration time.Duration) (bool, error) {
	var ok bool
	err := s.pg.QueryRow(ctx, `
INSERT INTO tfllab1.user_lock AS lock (user_id, instance_id, expires_at) VALUES ($1, $2, now()+$3)
    ON CONFLICT (user_id) DO UPDATE
        SET (expires_at, instance_id) = (excluded.expires_at, excluded.instance_id)
        WHERE expires_at > now()
    RETURNING true;
`,
		userID, instanceID, duration).Scan(&ok)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *PostgresStorage) Unlock(ctx context.Context, userID int64, instanceID uuid.UUID) error {
	_, err := s.pg.Exec(ctx, "DELETE FROM tfllab1.user_lock WHERE user_id = $1 AND instance_id = $2", userID, instanceID)
	return err
}

func (s *PostgresStorage) ForceUnlock(ctx context.Context, userID int64) error {
	_, err := s.pg.Exec(ctx, "DELETE FROM tfllab1.user_lock WHERE user_id = $1", userID)
	return err
}

func (s *PostgresStorage) IsLocked(ctx context.Context, userID int64, instanceID uuid.UUID) bool {
	var ok bool
	err := s.pg.QueryRow(ctx, "SELECT true FROM tfllab1.user_lock WHERE user_id = $1 AND instance_id = $2 AND expires_at > now()",
		userID, instanceID).Scan(&ok)
	if errors.Is(err, pgx.ErrNoRows) {
		return false
	}
	if err != nil {
		slog.Error("IsLocked request failed", "error", err)
		return false
	}
	return ok
}
