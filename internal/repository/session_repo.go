package repository

import (
	"context"
	"database/sql"

	u "github.com/cheezecakee/go-backend-utils/pkg/util"
	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
)

type SessionRepo interface {
	Create(ctx context.Context, session *m.Session) (*m.Session, error)
	GetByID(ctx context.Context, id uuid.UUID) (*m.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Session, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*m.Session, error)
	Update(ctx context.Context, session *m.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*m.Session, error)
	ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Session, error)
}

type DBSessionRepo struct {
	*u.BaseRepository
}

func NewSessionRepo(db *sql.DB) SessionRepo {
	return &DBSessionRepo{
		u.NewBaseRepository(db),
	}
}

const createSesssion = `INSERT INTO sessions (user_id, plan_id, name, start_time) VALUES ($1, $2, $3, NOW()) RETURNING id`

func (r *DBSessionRepo) Create(ctx context.Context, session *m.Session) (*m.Session, error) {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createSesssion, session.UserID, session.PlanID, session.Name).Scan(session.ID)
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}

const getSessionByID = `SELECT id, user_id, plan_id, name, start_time, end_time, notes, created_at, updated_at FROM sessions WHERE id = $1`

func (r *DBSessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*m.Session, error) {
	session := &m.Session{}
	row := r.DB.QueryRowContext(ctx, getSessionByID, id)
	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.PlanID,
		&session.Name,
		&session.StartTime,
		&session.EndTime,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	return session, err
}

const getSessionByUserID = `SELECT id, user_id, plan_id, name, start_time, end_time, notes, created_at, updated_at FROM sessions WHERE user_id = $1`

func (r *DBSessionRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Session, error) {
	rows, err := r.DB.QueryContext(ctx, getSessionByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*m.Session

	for rows.Next() {
		session := &m.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.PlanID,
			&session.Name,
			&session.StartTime,
			&session.EndTime,
			&session.Notes,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessions, err
}

const getSessionByPlanID = `SELECT id, user_id, plan_id, name, start_time, end_time, notes, created_at, updated_at FROM sessions WHERE plan_id = $1`

func (r *DBSessionRepo) GetByPlanID(ctx context.Context, planID uint) ([]*m.Session, error) {
	rows, err := r.DB.QueryContext(ctx, getSessionByPlanID, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*m.Session

	for rows.Next() {
		session := &m.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.PlanID,
			&session.Name,
			&session.StartTime,
			&session.EndTime,
			&session.Notes,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessions, err
}

const updateSession = `UPDATE sessions
SET
    notes = COALESCE(NULLIF($2, ''), notes),
    updated_at = NOW()
WHERE id = $1`

func (r *DBSessionRepo) Update(ctx context.Context, session *m.Session) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, updateSession, session.ID, session.Notes)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const deleteSession = `DELETE FROM sessions WHERE id = $1`

func (r *DBSessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteSession, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const listSession = `SELECT id, user_id, plan_id, name, start_time, end_time, notes, created_at, updated_at FROM sessions OFFSET $1 LIMIT $2`

func (r *DBSessionRepo) List(ctx context.Context, offset, limit int) ([]*m.Session, error) {
	rows, err := r.DB.QueryContext(ctx, listSession, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*m.Session

	for rows.Next() {
		session := &m.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.PlanID,
			&session.Name,
			&session.StartTime,
			&session.EndTime,
			&session.Notes,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessions, err
}

const listSessionByDateRange = `SELECT id, user_id, plan_id, name, start_time, end_time, notes, created_at, updated_at 
FROM sessions 
WHERE user_id = $1
AND created_at = BETWEEN $2 AND $3
ORDER BY created_at DESC`

func (r *DBSessionRepo) ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Session, error) {
	rows, err := r.DB.QueryContext(ctx, listSessionByDateRange, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*m.Session

	for rows.Next() {
		session := &m.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.PlanID,
			&session.Name,
			&session.StartTime,
			&session.EndTime,
			&session.Notes,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessions, err
}
