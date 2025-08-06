package exercise

import (
	"context"
	"database/sql"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type ExerciseRepo interface {
	Create(ctx context.Context, exercise *Exercise) error
	Update(ctx context.Context, exercise *Exercise) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*Exercise, error)
	Search(ctx context.Context, query string) ([]*Exercise, error)

	GetByID(ctx context.Context, id int) (*Exercise, error)
	GetByName(ctx context.Context, name string) (*Exercise, error)

	GetByCategoryID(ctx context.Context, category string) ([]*Exercise, error)
	GetByEquipmentName(ctx context.Context, equipment string) ([]*Exercise, error)

	// Exercise muscle groups
	AddMuscleGroups(ctx context.Context, exerciseID int, muscleGroupIDs []int) error
	RemoveMuscleGroups(ctx context.Context, exerciseID int, muscleGroupIDs []int) error
	GetMuscleGroups(ctx context.Context, exerciseID int) ([]*MuscleGroup, error)
	GetExercisesByMuscle(ctx context.Context, muscleGroupID int) ([]*Exercise, error)
	RemoveAllMuscleGroups(ctx context.Context, exerciseID int) error
	GetExercisesByMuscleName(ctx context.Context, muscleName string) ([]*Exercise, error)

	// Exercise training type groups
	AddExerciseTypes(ctx context.Context, exerciseID int, typeIDs []int) error
	RemoveExerciseTypes(ctx context.Context, exerciseID int, typeIDs []int) error
	GetExerciseTypes(ctx context.Context, exerciseID int) ([]*TrainingType, error)
	GetExercisesByType(ctx context.Context, typeID int) ([]*Exercise, error)
	RemoveAllExerciseTypes(ctx context.Context, exerciseID int) error
	GetExercisesByTypeName(ctx context.Context, typeName string) ([]*Exercise, error)
}

type exerciseRepo struct {
	tx transaction.BaseRepository
}

func NewExerciseRepo(db *sql.DB) ExerciseRepo {
	return &exerciseRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createExercise = `INSERT INTO exercises (name, description, category_id, equipment_id) VALUES ($1, $2, $3, $4) RETURNING id`

func (r *exerciseRepo) Create(ctx context.Context, exercise *Exercise) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createExercise, exercise.Name, exercise.Description, exercise.Category.ID, exercise.Equipment.ID).Scan(&exercise.ID)
	})
	return err
}

const updateExercise = `UPDATE exercises
SET 
    name = $2,
    description = $3,
    category_id = $4,
    equipment_id = $5,
    updated_at = NOW()
WHERE id = $1`

func (r *exerciseRepo) Update(ctx context.Context, exercise *Exercise) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, updateExercise, exercise.ID, exercise.Name, exercise.Description, exercise.Category.ID, exercise.Equipment.ID)
		return err
	})
	return err
}

const deleteExercise = `DELETE FROM exercises WHERE id = $1`

func (r *exerciseRepo) Delete(ctx context.Context, id int) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteExercise, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const listExercises = `SELECT id, name, description, category_id, equipment_id, created_at, updated_at FROM exercises OFFSET $1 LIMIT $2`

func (r *exerciseRepo) List(ctx context.Context, offset, limit int) ([]*Exercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, listExercises, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise
	for rows.Next() {
		exercise := &Exercise{
			Category:  &Category{},
			Equipment: &Equipment{},
		}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category.ID,
			&exercise.Equipment.ID,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, rows.Err()
}

const searchExercise = `
SELECT e.id, e.name, e.description, e.category_id, e.equipment_id, e.created_at, e.updated_at
FROM exercises e
JOIN exercise_categories c ON e.category_id = c.id
JOIN equipment eq ON e.equipment_id = eq.id
WHERE
    e.name ILIKE '%' || $1 || '%' OR
    e.description ILIKE '%' || $1 || '%' OR
    c.name ILIKE '%' || $1 || '%' OR
    eq.name ILIKE '%' || $1 || '%'`

func (r *exerciseRepo) Search(ctx context.Context, query string) ([]*Exercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, searchExercise, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise
	for rows.Next() {
		exercise := &Exercise{
			Category:  &Category{},
			Equipment: &Equipment{},
		}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category.ID,
			&exercise.Equipment.ID,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, rows.Err()
}

const getExerciseByID = `SELECT id, name, description, category_id, equipment_id, created_at, updated_at FROM exercises WHERE id = $1 LIMIT 1`

func (r *exerciseRepo) GetByID(ctx context.Context, id int) (*Exercise, error) {
	exercise := &Exercise{
		Category:  &Category{},
		Equipment: &Equipment{},
	}
	row := r.tx.DB().QueryRowContext(ctx, getExerciseByID, id)
	err := row.Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.Category.ID,
		&exercise.Equipment.ID,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	return exercise, err
}

const getExerciseByName = `SELECT id, name, description, category_id, equipment_id, created_at, updated_at FROM exercises WHERE name = $1 LIMIT 1`

func (r *exerciseRepo) GetByName(ctx context.Context, name string) (*Exercise, error) {
	exercise := &Exercise{
		Category:  &Category{},
		Equipment: &Equipment{},
	}
	row := r.tx.DB().QueryRowContext(ctx, getExerciseByName, name)
	err := row.Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.Category.ID,
		&exercise.Equipment.ID,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	return exercise, err
}

const getByCategoryID = `
SELECT e.id, e.name, e.description, e.category_id, e.equipment_id, e.created_at, e.updated_at 
FROM exercises e 
JOIN exercise_categories c ON e.category_id = c.id
WHERE c.name = $1`

func (r *exerciseRepo) GetByCategoryID(ctx context.Context, category string) ([]*Exercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getByCategoryID, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise
	for rows.Next() {
		exercise := &Exercise{
			Category:  &Category{},
			Equipment: &Equipment{},
		}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category.ID,
			&exercise.Equipment.ID,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, rows.Err()
}

const getByEquipmentName = `
SELECT e.id, e.name, e.description, e.category_id, e.equipment_id, e.created_at, e.updated_at 
FROM exercises e
JOIN equipment eq ON e.equipment_id = eq.id
WHERE eq.name = $1`

func (r *exerciseRepo) GetByEquipmentName(ctx context.Context, equipment string) ([]*Exercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getByEquipmentName, equipment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*Exercise
	for rows.Next() {
		exercise := &Exercise{
			Category:  &Category{},
			Equipment: &Equipment{},
		}
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category.ID,
			&exercise.Equipment.ID,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, rows.Err()
}
