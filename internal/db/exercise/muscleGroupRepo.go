package exercise

import (
	"context"
	"database/sql"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type MuscleGroupRepo interface {
	Create(ctx context.Context, muscleGroup *MuscleGroup) error
	GetByID(ctx context.Context, id int) (*MuscleGroup, error)
	GetByName(ctx context.Context, name string) (*MuscleGroup, error)
	List(ctx context.Context, offset, limit int) ([]*MuscleGroup, error)
}

type muscleGroupRepo struct {
	tx transaction.BaseRepository
}

func NewMuscleGroupRepo(db *sql.DB) MuscleGroupRepo {
	return &muscleGroupRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createMuscleGroup = `INSERT INTO muscle_groups (name) VALUES ($1) RETURNING id`

func (r *muscleGroupRepo) Create(ctx context.Context, muscleGroup *MuscleGroup) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createMuscleGroup, muscleGroup.Name).Scan(&muscleGroup.ID)
	})
}

const getMuscleGroupByID = `SELECT id, name FROM muscle_groups WHERE id = $1 LIMIT 1`

func (r *muscleGroupRepo) GetByID(ctx context.Context, id int) (*MuscleGroup, error) {
	muscleGroup := &MuscleGroup{}
	err := r.tx.DB().QueryRowContext(ctx, getMuscleGroupByID, id).Scan(&muscleGroup.ID, &muscleGroup.Name)
	return muscleGroup, err
}

const getMuscleGroupByName = `SELECT id, name FROM muscle_groups WHERE name = $1 LIMIT 1`

func (r *muscleGroupRepo) GetByName(ctx context.Context, name string) (*MuscleGroup, error) {
	muscleGroup := &MuscleGroup{}
	err := r.tx.DB().QueryRowContext(ctx, getMuscleGroupByName, name).Scan(&muscleGroup.ID, &muscleGroup.Name)
	return muscleGroup, err
}

const listMuscleGroups = `SELECT id, name FROM muscle_groups OFFSET $1 LIMIT $2`

func (r *muscleGroupRepo) List(ctx context.Context, offset, limit int) ([]*MuscleGroup, error) {
	rows, err := r.tx.DB().QueryContext(ctx, listMuscleGroups, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var muscleGroups []*MuscleGroup
	for rows.Next() {
		muscleGroup := &MuscleGroup{}
		if err := rows.Scan(&muscleGroup.ID, &muscleGroup.Name); err != nil {
			return nil, err
		}
		muscleGroups = append(muscleGroups, muscleGroup)
	}
	return muscleGroups, rows.Err()
}
