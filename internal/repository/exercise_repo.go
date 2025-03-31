package repository

import (
	"context"
	"database/sql"

	u "github.com/cheezecakee/go-backend-utils/pkg/util"

	m "github.com/cheezecakee/fitrkr/internal/models"
)

type ExerciseRepo interface {
	Create(ctx context.Context, exercise *m.Exercise) error
	GetByID(ctx context.Context, id uint) (*m.Exercise, error)
	GetByName(ctx context.Context, name string) (*m.Exercise, error)
	GetByCategory(ctx context.Context, category string) ([]*m.Exercise, error)
	GetByEquipment(ctx context.Context, equipment string) ([]*m.Exercise, error)
	Update(ctx context.Context, exercise *m.Exercise) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*m.Exercise, error)
	Search(ctx context.Context, query string) ([]*m.Exercise, error)
}

type DBExerciseRepo struct {
	*u.BaseRepository
}

func NewExerciseRepo(db *sql.DB) ExerciseRepo {
	return &DBExerciseRepo{
		u.NewBaseRepository(db),
	}
}

const createExercise = `INSERT INTO exercises (name, description, category, equipment) VALUES ($1, $2, $3, $4)`

func (r *DBExerciseRepo) Create(ctx context.Context, exercise *m.Exercise) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, createExercise, exercise.Name, exercise.Description, exercise.Category, exercise.Equipment)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const getExByExerciseID = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises WHERE id = $1 LIMIT 1`

func (r *DBExerciseRepo) GetByID(ctx context.Context, id uint) (*m.Exercise, error) {
	exercise := &m.Exercise{}
	row := r.DB.QueryRowContext(ctx, getExByExerciseID, id)
	err := row.Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.Category,
		&exercise.Equipment,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	return exercise, err
}

const getExByExerciseName = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises WHERE name = $1 LIMIT 1`

func (r *DBExerciseRepo) GetByName(ctx context.Context, name string) (*m.Exercise, error) {
	exercise := &m.Exercise{}
	row := r.DB.QueryRowContext(ctx, getExByExerciseName, name)
	err := row.Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.Category,
		&exercise.Equipment,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	return exercise, err
}

const getExByExerciseCategory = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises WHERE category = $1`

func (r *DBExerciseRepo) GetByCategory(ctx context.Context, category string) ([]*m.Exercise, error) {
	rows, err := r.DB.QueryContext(ctx, getExByExerciseCategory, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*m.Exercise
	for rows.Next() {
		exercise := &m.Exercise{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category,
			&exercise.Equipment,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exercises, err
}

const getExByExerciseEquipment = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises WHERE equipment = $1`

func (r *DBExerciseRepo) GetByEquipment(ctx context.Context, equipment string) ([]*m.Exercise, error) {
	rows, err := r.DB.QueryContext(ctx, getExByExerciseEquipment, equipment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*m.Exercise
	for rows.Next() {
		exercise := &m.Exercise{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category,
			&exercise.Equipment,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exercises, err
}

const updateExercise = `UPDATE exercises
SET 
    name = $2,
    description = $3,
    category = $4,
    equipment = $5,
    updated_at = NOW()
WHERE id = $1`

func (r *DBExerciseRepo) Update(ctx context.Context, exercise *m.Exercise) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, updateExercise, exercise.ID, exercise.Name, exercise.Description, exercise.Category, exercise.Equipment)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const deleteExercise = `DELETE FROM exercises WHERE id = $1`

func (r *DBExerciseRepo) Delete(ctx context.Context, id uint) error {
	err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteExercise, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const listExercises = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises OFFSET $1 LIMIT $2`

func (r *DBExerciseRepo) List(ctx context.Context, offset, limit int) ([]*m.Exercise, error) {
	rows, err := r.DB.QueryContext(ctx, listExercises, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*m.Exercise
	for rows.Next() {
		exercise := &m.Exercise{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category,
			&exercise.Equipment,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exercises, err
}

const searchExercise = `SELECT id, name, description, category, equipment, created_at, updated_at
FROM exercises
WHERE
    name ILIKE '%' || $1 || '%' OR
    description ILIKE '%' || $1 || '%' OR
    category ILIKE '%' || $1 || '%' OR
    equipment ILIKE '%' || $1 || '%'`

func (r *DBExerciseRepo) Search(ctx context.Context, query string) ([]*m.Exercise, error) {
	rows, err := r.DB.QueryContext(ctx, searchExercise, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*m.Exercise
	for rows.Next() {
		exercise := &m.Exercise{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category,
			&exercise.Equipment,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exercises, err
}
