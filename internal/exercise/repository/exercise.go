package repository

import (
	"context"
	"database/sql"

	m "github/cheezecakee/fitrkr/internal/exercise/models"
)

type DBExerciseRepo struct {
	db *sql.DB
}

func New(db *sql.DB) ExerciseRepo {
	return &DBExerciseRepo{
		db: db,
	}
}

const createExercise = `INSERT INTO exercises (name, description, category, equipment) VALUES ($1, $2, $3, $4)`

func (r *DBExerciseRepo) Create(ctx context.Context, exercise *m.Exercise) error {
	_, err := r.db.ExecContext(ctx, createExercise, exercise.Name, exercise.Description, exercise.Category, exercise.Equipment)
	return err
}

const getByExerciseID = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises WHERE id = $1 LIMIT 1`

func (r *DBExerciseRepo) GetByID(ctx context.Context, id uint) (*m.Exercise, error) {
	exercise := &m.Exercise{}
	row := r.db.QueryRowContext(ctx, getByExerciseID, id)
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

const getByExerciseName = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises WHERE name = $1 LIMIT 1`

func (r *DBExerciseRepo) GetByName(ctx context.Context, name string) (*m.Exercise, error) {
	exercise := &m.Exercise{}
	row := r.db.QueryRowContext(ctx, getByExerciseName, name)
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

const getByExerciseCategory = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises WHERE category = $1`

func (r *DBExerciseRepo) GetByCategory(ctx context.Context, category string) ([]*m.Exercise, error) {
	rows, err := r.db.QueryContext(ctx, getByExerciseCategory, category)
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

const getByExerciseEquipment = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises WHERE equipment = $1`

func (r *DBExerciseRepo) GetByEquipment(ctx context.Context, equipment string) ([]*m.Exercise, error) {
	rows, err := r.db.QueryContext(ctx, getByExerciseEquipment, equipment)
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
	_, err := r.db.ExecContext(ctx, updateExercise, exercise.ID, exercise.Name, exercise.Description, exercise.Category, exercise.Equipment)
	return err
}

const deleteExercise = `DELETE FROM exercises WHERE id = $1`

func (r *DBExerciseRepo) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, deleteExercise, id)
	return err
}

const listExercises = `SELECT id, name, description, category, equipment,  created_at, updated_at FROM exercises OFFSET $1 LIMIT $2`

func (r *DBExerciseRepo) List(ctx context.Context, offset, limit int) ([]*m.Exercise, error) {
	rows, err := r.db.QueryContext(ctx, listExercises, offset, limit)
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
	rows, err := r.db.QueryContext(ctx, searchExercise, query)
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
