package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/session/models"
)

type DBSessionRepo struct {
	db *sql.DB
}

func New(db *sql.DB) SessionRepo {
	return &DBSessionRepo{
		db: db,
	}
}

const createSesssion = `INSERT INTO sessions (user_id, plan_id, name, start_time) VALUES ($1, $2, $3, NOW()) RETURNING id`

func (r *DBSessionRepo) Create(ctx context.Context, session *m.Session) (*m.Session, error) {
	err := r.db.QueryRowContext(ctx, createSesssion, session.UserID, session.PlanID, session.Name).Scan(session.ID)
	if err != nil {
		return nil, err
	}
	return session, nil
}

const getSessionByID = `SELECT id, user_id, plan_id, name, start_time, end_time, notes, created_at, updated_at FROM sessions WHERE id = $1`

func (r *DBSessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*m.Session, error) {
	session := &m.Session{}
	row := r.db.QueryRowContext(ctx, getSessionByID, id)
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
	rows, err := r.db.QueryContext(ctx, getSessionByUserID, userID)
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
	rows, err := r.db.QueryContext(ctx, getSessionByPlanID, planID)
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

const updateSession = ``

func (r *DBSessionRepo) Update(ctx context.Context, session *m.Session) error {}

const deleteSession = `DELETE FROM sessions WHERE id = $1`

func (r *DBSessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteSession, id)
	return err
}

const listSession = `SELECT id, user_id, plan_id, name, start_time, end_time, notes, created_at, updated_at FROM sessions OFFSET $1 LIMIT $2`

func (r *DBSessionRepo) List(ctx context.Context, offset, limit int) ([]*m.Session, error) {
	rows, err := r.db.QueryContext(ctx, listSession, offset, limit)
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
	rows, err := r.db.QueryContext(ctx, listSessionByDateRange, userID, startDate, endDate)
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
