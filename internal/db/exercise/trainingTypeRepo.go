package exercise

import (
	"context"
	"database/sql"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type TrainingTypeRepo interface {
	Create(ctx context.Context, exerciseType *TrainingType) error
	GetByName(ctx context.Context, name string) (*TrainingType, error)
	GetByID(ctx context.Context, id int) (*TrainingType, error)
	List(ctx context.Context, offset, limit int) ([]*TrainingType, error)
}

type trainingTypeRepo struct {
	tx transaction.BaseRepository
}

func NewTrainingTypeRepo(db *sql.DB) TrainingTypeRepo {
	return &trainingTypeRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createTrainingType = `INSERT INTO training_types (name) VALUES ($1) RETURNING id`

func (r *trainingTypeRepo) Create(ctx context.Context, exerciseType *TrainingType) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createTrainingType, exerciseType.Name).Scan(&exerciseType.ID)
	})
}

const getTrainingTypeByID = `SELECT id, name FROM training_types WHERE id = $1 LIMIT 1`

func (r *trainingTypeRepo) GetByID(ctx context.Context, id int) (*TrainingType, error) {
	exerciseType := &TrainingType{}
	err := r.tx.DB().QueryRowContext(ctx, getTrainingTypeByID, id).Scan(&exerciseType.ID, &exerciseType.Name)
	return exerciseType, err
}

const getTrainingTypeByName = `SELECT id, name FROM training_types WHERE name = $1 LIMIT 1`

func (r *trainingTypeRepo) GetByName(ctx context.Context, name string) (*TrainingType, error) {
	exerciseType := &TrainingType{}
	err := r.tx.DB().QueryRowContext(ctx, getTrainingTypeByName, name).Scan(&exerciseType.ID, &exerciseType.Name)
	return exerciseType, err
}

const listTrainingTypes = `SELECT id, name FROM training_types OFFSET $1 LIMIT $2`

func (r *trainingTypeRepo) List(ctx context.Context, offset, limit int) ([]*TrainingType, error) {
	rows, err := r.tx.DB().QueryContext(ctx, listTrainingTypes, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exerciseTypes []*TrainingType
	for rows.Next() {
		exerciseType := &TrainingType{}
		if err := rows.Scan(&exerciseType.ID, &exerciseType.Name); err != nil {
			return nil, err
		}
		exerciseTypes = append(exerciseTypes, exerciseType)
	}
	return exerciseTypes, rows.Err()
}
