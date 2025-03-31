package repository

import (
	"context"
	"database/sql"

	u "github.com/cheezecakee/go-backend-utils/pkg/util"
	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
)

type PlanRepo interface {
	Create(ctx context.Context, plan *m.Plan) (*m.Plan, error)
	GetByID(ctx context.Context, id uint) (*m.Plan, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Plan, error)
	Update(ctx context.Context, plan *m.Plan) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*m.Plan, error)
}

type DBPlanRepo struct {
	*u.BaseRepository
}

func NewPlanRepo(db *sql.DB) PlanRepo {
	return &DBPlanRepo{
		u.NewBaseRepository(db),
	}
}

const createPlan = `INSERT INTO plans (
    user_id, 
    name, 
    description, 
    created_at, 
    updated_at
)
VALUES (
    $1, --user_id
    $2, --name
    NULLIF($3,  '') --description
    ),
    NOW(), --created_at
    NOW()  --updated_at
)
RETURNING id`

func (r *DBPlanRepo) Create(ctx context.Context, plan *m.Plan) (*m.Plan, error) {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createPlan, plan.UserID, plan.Name, plan.Description).Scan(&plan.ID)
	})
	if err != nil {
		return nil, err
	}
	return plan, nil
}

const getPlanByID = `SELECT id, user_id, name, description, is_active, created_at, updated_at FROM plans WHERE id = $1`

func (r *DBPlanRepo) GetByID(ctx context.Context, id uint) (*m.Plan, error) {
	plan := &m.Plan{}
	row := r.DB.QueryRowContext(ctx, getPlanByID, id)
	err := row.Scan(
		&plan.ID,
		&plan.UserID,
		&plan.Name,
		&plan.Description,
		&plan.IsActive,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	return plan, err
}

const getPlanByUserID = `SELECT id, user_id, name, description, is_active, created_at, updated_at FROM plans WHERE user_id = $1`

func (r *DBPlanRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Plan, error) {
	rows, err := r.DB.QueryContext(ctx, getPlanByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*m.Plan
	for rows.Next() {
		plan := &m.Plan{}
		err := rows.Scan(
			&plan.ID,
			&plan.UserID,
			&plan.Name,
			&plan.Description,
			&plan.IsActive,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

const updatePlan = `UPDATE plans
SET name = COALESCE(NULLIF($3,  ''), name),
    description = COALESCE(NULLIF($4,  ''), description,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2`

func (r *DBPlanRepo) Update(ctx context.Context, plan *m.Plan) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, updatePlan, plan.ID, plan.UserID, plan.Name, plan.Description)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const deletePlan = `DELETE FROM plans WHERE id = $1`

func (r *DBPlanRepo) Delete(ctx context.Context, id uint) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deletePlan, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const listPlan = `SELECT id, user_id, name, description, is_active, created_at, updated_at 
FROM plans 
ORDER BY created_at DESC
OFFSET $1 LIMIT $2 `

func (r *DBPlanRepo) List(ctx context.Context, offset, limit int) ([]*m.Plan, error) {
	rows, err := r.DB.QueryContext(ctx, listPlan, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*m.Plan
	for rows.Next() {
		plan := &m.Plan{}
		err := rows.Scan(
			&plan.ID,
			&plan.UserID,
			&plan.Name,
			&plan.Description,
			&plan.IsActive,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}
