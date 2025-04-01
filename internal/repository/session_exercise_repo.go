package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type SessionExRepo interface {
	Create(ctx context.Context, sessionEx *m.SessionEx) (*m.SessionEx, error)
	GetByID(ctx context.Context, id uuid.UUID) (*m.SessionEx, error)
	GetBysessionID(ctx context.Context, sessionID uuid.UUID) ([]*m.SessionEx, error)
	Update(ctx context.Context, sessionEx *m.SessionEx) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBysessionID(ctx context.Context, sessionID uuid.UUID) error
}

type DBSessionExRepo struct {
	tx transaction.BaseRepository
}

func NewSessionExRepo(db *sql.DB) SessionExRepo {
	return &DBSessionExRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createSessionEx = `INSERT INTO sessions_exercises (session_id, exercise_id, "order") VALUES ($1, $2, $3) RETURNING id`

func (r *DBSessionExRepo) Create(ctx context.Context, sessionEx *m.SessionEx) (*m.SessionEx, error) {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createSessionEx, sessionEx.SessionID, sessionEx.ExerciseID, sessionEx.Order).Scan(sessionEx.ID)
	})
	if err != nil {
		return nil, err
	}

	return sessionEx, nil
}

const getSessionExByID = `SELECT id, session_id, exercise_id, "order", created_at, updated_at FROM sessions_exercises WHERE id = $1 LIMIT 1`

func (r *DBSessionExRepo) GetByID(ctx context.Context, id uuid.UUID) (*m.SessionEx, error) {
	row := r.tx.DB().QueryRowContext(ctx, getSessionExByID, id)
	sessionEx := &m.SessionEx{}
	err := row.Scan(
		&sessionEx.ID,
		&sessionEx.SessionID,
		&sessionEx.ExerciseID,
		&sessionEx.Order,
		&sessionEx.CreatedAt,
		&sessionEx.UpdatedAt,
	)

	return sessionEx, err
}

const getSessionExBySessionID = `SELECT id, session_id, exercise_id, "order", created_at, updated_at FROM sessions_exercises WHERE session_id = $1`

func (r *DBSessionExRepo) GetBysessionID(ctx context.Context, sessionID uuid.UUID) ([]*m.SessionEx, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getSessionExBySessionID, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessionExs []*m.SessionEx

	for rows.Next() {
		sessionEx := &m.SessionEx{}
		err := rows.Scan(
			&sessionEx.ID,
			&sessionEx.SessionID,
			&sessionEx.ExerciseID,
			&sessionEx.Order,
			&sessionEx.CreatedAt,
			&sessionEx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sessionExs = append(sessionExs, sessionEx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessionExs, err
}

const updateSessionEx = `UPDATE sessions_exercises
SET "order" = $2, updated_at = NOW()
WHERE id = $1`

func (r *DBSessionExRepo) Update(ctx context.Context, sessionEx *m.SessionEx) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, updateSessionEx, sessionEx.ID, sessionEx.Order)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const deleteSessionEx = `DELETE FROM sessions_exercises WHERE id = $1`

func (r *DBSessionExRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteSessionEx, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const deleteSessionExBySessionID = `DELETE FROM sessions_exercises WHERE session_id = $1`

func (r *DBSessionExRepo) DeleteBysessionID(ctx context.Context, sessionID uuid.UUID) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteSessionEx, sessionID)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
