package repository

import (
	"context"
	"database/sql"

	u "github.com/cheezecakee/go-backend-utils/pkg/util"
	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
)

type LogRepo interface {
	Create(ctx context.Context, log *m.Log) error
	GetByID(ctx context.Context, id uint) (*m.Log, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Log, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*m.Log, error)
	GetByType(ctx context.Context, logType string) ([]*m.Log, error)
	GetByPriority(ctx context.Context, priority string) ([]*m.Log, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*m.Log, error)
	ListByRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Log, error)
	GetPRs(ctx context.Context, userID uuid.UUID) ([]*m.Log, error)
}

type DBLogRepo struct {
	*u.BaseRepository
}

func NewLogRepo(db *sql.DB) LogRepo {
	return &DBLogRepo{
		u.NewBaseRepository(db),
	}
}

const createLog = `INSERT INTO logs (user_id, plan_id, metadata, type, priority, message, pr) Values ($1, $2, $3, $4, $5, $6, $7)`

func (r *DBLogRepo) Create(ctx context.Context, log *m.Log) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, createLog, log.UserID, log.PlanID, log.Context, log.Type, log.Priority, log.Message, log.Pr)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const getLogByID = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE id = $1`

func (r *DBLogRepo) GetByID(ctx context.Context, id uint) (*m.Log, error) {
	log := &m.Log{}
	row := r.DB.QueryRowContext(ctx, getLogByID, id)
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

const getLogByUserID = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE user_id = $1`

func (r *DBLogRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Log, error) {
	rows, err := r.DB.QueryContext(ctx, getLogByUserID, userID)
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

const getLogByPlanID = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE plan_id = $1`

func (r *DBLogRepo) GetByPlanID(ctx context.Context, planID uint) ([]*m.Log, error) {
	rows, err := r.DB.QueryContext(ctx, getLogByPlanID, planID)
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

const getLogByType = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE type = $1`

func (r *DBLogRepo) GetByType(ctx context.Context, logType string) ([]*m.Log, error) {
	rows, err := r.DB.QueryContext(ctx, getLogByType, logType)
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

const getLogByPriority = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs WHERE priority = $1`

func (r *DBLogRepo) GetByPriority(ctx context.Context, priority string) ([]*m.Log, error) {
	rows, err := r.DB.QueryContext(ctx, getLogByPriority, priority)
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
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteLog, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const listLog = `SELECT id, user_id, plan_id, metadata, type, priority, message, pr FROM logs OFFSET $1 LIMIT $2`

func (r *DBLogRepo) List(ctx context.Context, offset, limit int) ([]*m.Log, error) {
	rows, err := r.DB.QueryContext(ctx, listLog, offset, limit)
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
	rows, err := r.DB.QueryContext(ctx, listByLogRange, userID, startDate, endDate)
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

const getLogPRs = `SELECT pr FROM logs WHERE user_id = $1`

func (r *DBLogRepo) GetPRs(ctx context.Context, userID uuid.UUID) ([]*m.Log, error) {
	rows, err := r.DB.QueryContext(ctx, getLogPRs, userID)
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
