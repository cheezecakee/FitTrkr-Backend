package exercise

import (
	"context"
	"database/sql"
)

const addExerciseTrainingTypes = `INSERT INTO exercise_training_types (exercise_id, training_type_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

func (r *exerciseRepo) AddExerciseTypes(ctx context.Context, exerciseID int, typeIDs []int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		for _, typeID := range typeIDs {
			_, err := tx.ExecContext(ctx, addExerciseTrainingTypes, exerciseID, typeID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

const removeExerciseTrainingTypes = `DELETE FROM exercise_training_types WHERE exercise_id = $1 AND training_type_id = $2`

func (r *exerciseRepo) RemoveExerciseTypes(ctx context.Context, exerciseID int, typeIDs []int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		for _, typeID := range typeIDs {
			_, err := tx.ExecContext(ctx, removeExerciseTrainingTypes, exerciseID, typeID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

const removeAllExerciseTrainingTypes = `DELETE FROM exercise_training_types WHERE exercise_id = $1`

func (r *exerciseRepo) RemoveAllExerciseTypes(ctx context.Context, exerciseID int) error {
	return r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, removeAllExerciseTrainingTypes, exerciseID)
		return err
	})
}

const getTrainingTypesByExercise = `
SELECT tt.id, tt.name
FROM training_types tt
JOIN exercise_training_types ett ON tt.id = ett.training_type_id
WHERE ett.exercise_id = $1
ORDER BY tt.name`

func (r *exerciseRepo) GetExerciseTypes(ctx context.Context, exerciseID int) ([]*TrainingType, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getTrainingTypesByExercise, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainingTypes []*TrainingType
	for rows.Next() {
		trainingType := &TrainingType{}
		if err := rows.Scan(&trainingType.ID, &trainingType.Name); err != nil {
			return nil, err
		}
		trainingTypes = append(trainingTypes, trainingType)
	}
	return trainingTypes, rows.Err()
}

const getExercisesByTrainingType = `
SELECT e.id, e.name, e.description, e.category_id, e.equipment_id, e.created_at, e.updated_at
FROM exercises e
JOIN exercise_training_types ett ON e.id = ett.exercise_id
WHERE ett.training_type_id = $1
ORDER BY e.name`

func (r *exerciseRepo) GetExercisesByType(ctx context.Context, typeID int) ([]*Exercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getExercisesByTrainingType, typeID)
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

const getExercisesByTrainingTypeName = `
SELECT e.id, e.name, e.description, e.category_id, e.equipment_id, e.created_at, e.updated_at
FROM exercises e
JOIN exercise_training_types ett ON e.id = ett.exercise_id
JOIN training_types tt ON ett.training_type_id = tt.id
WHERE tt.name = $1
ORDER BY e.name`

func (r *exerciseRepo) GetExercisesByTypeName(ctx context.Context, typeName string) ([]*Exercise, error) {
	rows, err := r.tx.DB().QueryContext(ctx, getExercisesByTrainingTypeName, typeName)
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
