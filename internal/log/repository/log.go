package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/log/models"
)

type DBLogRepo struct {
	db *sql.DB
}

func New(db *sql.DB) LogRepo {
	return &DBLogRepo{
		db: db,
	}
}

const createLog = `INSERT INTO logs (user_id, plan_id, metadata, type, priority, message, pr) Values ($1, $2, $3, $4, $5, $6, $7)`

func (r *DBLogRepo) Create(ctx context.Context, log *m.Log) error {
	_, err := r.db.ExecContext(ctx, createLog, log.UserID, log.PlanID, log.Context, log.Type, log.Priority, log.Message, log.Pr)
	return err
}

const getByID = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE id = $1`

func (r *DBLogRepo) GetByID(ctx context.Context, id uint) (*m.Log, error) {
	log := &m.Log{}
	row := r.db.QueryRowContext(ctx, getByID, id)
	err := row.Scan(
		&log.ID,
		&log.UserID,
		&log.PlanID,
		&log.Context,
		&log.Type,
		&log.Priority,
		&log.Message,
		&log.Pr,
	)
	return log, err
}

const getByUserID = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE user_id = $1`

func (r *DBLogRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Log, error) {
	rows, err := r.db.QueryContext(ctx, getByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*m.Log
	for rows.Next() {
		log := &m.Log{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.PlanID,
			&log.Context,
			&log.Type,
			&log.Priority,
			&log.Message,
			&log.Pr,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, err
}

const getByPlanID = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE plan_id = $1`

func (r *DBLogRepo) GetByPlanID(ctx context.Context, planID uint) ([]*m.Log, error) {
	rows, err := r.db.QueryContext(ctx, getByPlanID, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*m.Log
	for rows.Next() {
		log := &m.Log{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.PlanID,
			&log.Context,
			&log.Type,
			&log.Priority,
			&log.Message,
			&log.Pr,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, err
}

const getByType = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE type = $1`

func (r *DBLogRepo) GetByType(ctx context.Context, logType string) ([]*m.Log, error) {
	rows, err := r.db.QueryContext(ctx, getByType, logType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*m.Log
	for rows.Next() {
		log := &m.Log{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.PlanID,
			&log.Context,
			&log.Type,
			&log.Priority,
			&log.Message,
			&log.Pr,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, err
}

const getByPriority = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE priority = $1`

func (r *DBLogRepo) GetByPriority(ctx context.Context, priority string) ([]*m.Log, error) {
	rows, err := r.db.QueryContext(ctx, getByPriority, priority)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*m.Log
	for rows.Next() {
		log := &m.Log{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.PlanID,
			&log.Context,
			&log.Type,
			&log.Priority,
			&log.Message,
			&log.Pr,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, err
}

const deleteLog = `DELETE FROM logs WHERE id = $1`

func (r *DBLogRepo) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, deleteLog, id)
	return err
}

const listLog = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs OFFSET $1 LIMIT $2`

func (r *DBLogRepo) List(ctx context.Context, offset, limit int) ([]*m.Log, error) {
	rows, err := r.db.QueryContext(ctx, listLog, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*m.Log
	for rows.Next() {
		log := &m.Log{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.PlanID,
			&log.Context,
			&log.Type,
			&log.Priority,
			&log.Message,
			&log.Pr,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, err
}

const listByLogRange = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr 
FROM logs
WHERE user_id = $1
AND created_at BETWEEN $2 AND $3
ORDER BY created_at DESC`

func (r *DBLogRepo) ListByRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Log, error) {
	rows, err := r.db.QueryContext(ctx, listByLogRange, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*m.Log
	for rows.Next() {
		log := &m.Log{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.PlanID,
			&log.Context,
			&log.Type,
			&log.Priority,
			&log.Message,
			&log.Pr,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, err
}

const getPRs = `SELECT pr FROM logs WHERE user_id = $1`

func (r *DBLogRepo) GetPRs(ctx context.Context, userID uuid.UUID) ([]*m.Log, error) {
	rows, err := r.db.QueryContext(ctx, getPRs, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*m.Log
	for rows.Next() {
		log := &m.Log{}
		err := rows.Scan(
			&log.Pr,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, err
}
