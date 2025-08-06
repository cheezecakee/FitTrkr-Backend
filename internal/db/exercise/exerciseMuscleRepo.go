package exercise

import (
	"context"
	"database/sql"
)

const addExerciseMuscles = `INSERT INTO exercise_muscles (exercise_id, muscle_group_id) VALUES ($1, $2)`

func (r *exerciseRepo) AddMuscleGroups(ctx context.Context, exerciseID int, muscleGroupIDs []int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		for _, muscleGroupID := range muscleGroupIDs {
			_, err := tx.ExecContext(ctx, addExerciseMuscles, exerciseID, muscleGroupID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

const removeExerciseMuscles = `DELETE FROM exercise_muscles WHERE exercise_id = $1 AND muscle_group_id = $2`

func (r *exerciseRepo) RemoveMuscleGroups(ctx context.Context, exerciseID int, muscleGroupIDs []int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		for _, muscleGroupID := range muscleGroupIDs {
			_, err := tx.ExecContext(ctx, removeExerciseMuscles, exerciseID, muscleGroupID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

const removeAllExerciseMuscles = `DELETE FROM exercise_muscles WHERE exercise_id = $1`

func (r *exerciseRepo) RemoveAllMuscleGroups(ctx context.Context, exerciseID int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, removeAllExerciseMuscles, exerciseID)
		return err
	})
}

const getExerciseMuscles = `
SELECT m.id, m.name
FROM muscle_groups m
JOIN exercise_muscles em ON m.id = em.muscle_group_id
WHERE em.exercise_id = $1
ORDER BY m.name`

func (r *exerciseRepo) GetMuscleGroups(ctx context.Context, exerciseID int) ([]*MuscleGroup, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getExerciseMuscles, exerciseID)
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

const getExercisesByMuscleGroup = `
SELECT e.id, e.name, e.description, e.category_id, e.equipment_id, e.created_at, e.updated_at
FROM exercises e
JOIN exercise_muscles em ON e.id = em.exercise_id
WHERE em.muscle_group_id = $1
ORDER BY e.name`

func (r *exerciseRepo) GetExercisesByMuscle(ctx context.Context, muscleGroupID int) ([]*Exercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getExercisesByMuscleGroup, muscleGroupID)
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
		if err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category.ID,
			&exercise.Equipment.ID,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, rows.Err()
}

// FIXED: This should query muscle groups, not training types!
const getExercisesByMuscleName = `
SELECT e.id, e.name, e.description, e.category_id, e.equipment_id, e.created_at, e.updated_at
FROM exercises e
JOIN exercise_muscles em ON e.id = em.exercise_id
JOIN muscle_groups mg ON em.muscle_group_id = mg.id
WHERE mg.name = $1
ORDER BY e.name`

func (r *exerciseRepo) GetExercisesByMuscleName(ctx context.Context, muscleName string) ([]*Exercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getExercisesByMuscleName, muscleName)
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
		if err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category.ID,
			&exercise.Equipment.ID,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, rows.Err()
}
