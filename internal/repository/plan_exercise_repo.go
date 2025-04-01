package repository

import (
	"context"
	"database/sql"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type PlanExRepo interface {
	Create(ctx context.Context, planEx *m.PlanEx) (*m.PlanEx, error)
	GetByID(ctx context.Context, id uint) (*m.PlanEx, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*m.PlanEx, error)
	Update(ctx context.Context, planEx *m.PlanEx) error
	Delete(ctx context.Context, id uint) error
	DeleteByPlanID(ctx context.Context, planID uint) error
}

type DBPlanExRepo struct {
	tx transaction.BaseRepository
}

func NewPlanExRepo(db *sql.DB) PlanExRepo {
	return &DBPlanExRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createPlanEx = `
INSERT INTO plans_exercises (
    plan_id, 
    exercise_id, 
    name, 
    sets, 
    reps, 
    weight, 
    interval, 
    rest
)
VALUES (
    $1, -- plan_id
    $2, -- exercise_id
    $3, -- name
    $4, -- sets
    $5, -- reps
    $6, -- weight
    NULLIF($7,  ''), -- interval(optional)
    NULLIF($8,  '')   -- rest(optional)
) 
RETURNING id`

func (r *DBPlanExRepo) Create(ctx context.Context, planEx *m.PlanEx) (*m.PlanEx, error) {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createPlanEx, planEx.PlanID, planEx.ExerciseID, planEx.Name, planEx.Sets, planEx.Reps, planEx.Weight, planEx.Interval, planEx.Rest).Scan(&planEx.ID)
	})
	if err != nil {
		return nil, err
	}
	return planEx, nil
}

const getPlanExByID = `SELECT id, plan_id, exercise_id, name, sets, reps, weight, interval, rest, created_at, updated_at FROM plans_exercises WHERE id = $1`

func (r *DBPlanExRepo) GetByID(ctx context.Context, id uint) (*m.PlanEx, error) {
	planEx := &m.PlanEx{}
	row := r.tx.DB().QueryRowContext(ctx, getPlanExByID, id)
	err := row.Scan(
		&planEx.ID,
		&planEx.PlanID,
		&planEx.ExerciseID,
		&planEx.Name,
		&planEx.Sets,
		&planEx.Reps,
		&planEx.Weight,
		&planEx.Interval,
		&planEx.Rest,
		&planEx.CreatedAt,
		&planEx.UpdatedAt,
	)

	return planEx, err
}

const getPlanExByPlanID = `SELECT id, plan_id, exercise_id, name, sets, reps, weight, interval, rest, created_at, updated_at FROM plans_exercises WHERE plan_id = $1`

func (r *DBPlanExRepo) GetByPlanID(ctx context.Context, planID uint) ([]*m.PlanEx, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getPlanExByPlanID, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var planExs []*m.PlanEx
	for rows.Next() {
		planEx := &m.PlanEx{}
		err := rows.Scan(
			&planEx.ID,
			&planEx.PlanID,
			&planEx.ExerciseID,
			&planEx.Name,
			&planEx.Sets,
			&planEx.Reps,
			&planEx.Weight,
			&planEx.Interval,
			&planEx.Rest,
			&planEx.CreatedAt,
			&planEx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		planExs = append(planExs, planEx)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return planExs, err
}

const updatePlanEx = `UPDATE plans_exercises
SET exercise_id =  COALESCE(NULLIF($3,  ''), exercise_id),
    name = COALESCE(NULLIF($4,  ''), name),
    sets = COALESCE(NULLIF($5,  ''), sets),
    reps = COALESCE(NULLIF($6,  ''), reps),
    weight = COALESCE(NULLIF($7,  ''), weight),
    interval = COALESCE(NULLIF($8,  ''), interval),
    rest = COALESCE(NULLIF($9,  ''), rest),
    updated_at = NOW()
WHERE id = $1 AND plan_id = $2`

func (r *DBPlanExRepo) Update(ctx context.Context, planEx *m.PlanEx) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, updatePlanEx, planEx.ID, planEx.PlanID, planEx.ExerciseID, planEx.Name, planEx.Sets, planEx.Reps, planEx.Weight, planEx.Interval, planEx.Rest)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const deletePlanEx = `DELETE FROM plans_exercises WHERE id = $1`

func (r *DBPlanExRepo) Delete(ctx context.Context, id uint) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deletePlanEx, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const deletePlanExByPlanID = `DELETE FROM plans_exercises WHERE plan_id = $1`

func (r *DBPlanExRepo) DeleteByPlanID(ctx context.Context, planID uint) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deletePlanEx, planID)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
