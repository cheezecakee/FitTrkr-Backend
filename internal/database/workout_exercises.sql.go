// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: workout_exercises.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createWorkoutExercise = `-- name: CreateWorkoutExercise :one
INSERT INTO workout_exercises (
    workout_id,
    exercise_id,
    sets,
    reps_min,
    reps_max,
    weight,
    interval,
    rest
)

VALUES (
    $1, -- workout_id
    $2, -- exercise_id
    $3, -- sets
    $4, -- reps_min
    $5, -- reps_max
    $6, -- weight
    $7, -- interval
    $8  -- rest
)
RETURNING id, workout_id, exercise_id, sets, reps_min, reps_max, weight, interval, rest, created_at, updated_at
`

type CreateWorkoutExerciseParams struct {
	WorkoutID  uuid.UUID
	ExerciseID uuid.UUID
	Sets       int32
	RepsMin    int32
	RepsMax    int32
	Weight     float64
	Interval   int32
	Rest       int32
}

func (q *Queries) CreateWorkoutExercise(ctx context.Context, arg CreateWorkoutExerciseParams) (WorkoutExercise, error) {
	row := q.db.QueryRowContext(ctx, createWorkoutExercise,
		arg.WorkoutID,
		arg.ExerciseID,
		arg.Sets,
		arg.RepsMin,
		arg.RepsMax,
		arg.Weight,
		arg.Interval,
		arg.Rest,
	)
	var i WorkoutExercise
	err := row.Scan(
		&i.ID,
		&i.WorkoutID,
		&i.ExerciseID,
		&i.Sets,
		&i.RepsMin,
		&i.RepsMax,
		&i.Weight,
		&i.Interval,
		&i.Rest,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteWorkoutExercise = `-- name: DeleteWorkoutExercise :exec
DELETE FROM workout_exercises WHERE id = $1
`

func (q *Queries) DeleteWorkoutExercise(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteWorkoutExercise, id)
	return err
}

const editWorkoutExercise = `-- name: EditWorkoutExercise :one
UPDATE workout_exercises
SET sets = COALESCE($2, sets),
    reps_min = COALESCE($3, reps_min),
    reps_max = COALESCE($4, reps_max),
    weight = COALESCE($5, weight),
    interval = COALESCE($6, interval),
    rest = COALESCE($7, rest),
    updated_at = NOW()
WHERE id = $1
RETURNING id, workout_id, exercise_id, sets, reps_min, reps_max, weight, interval, rest, created_at, updated_at
`

type EditWorkoutExerciseParams struct {
	ID       uuid.UUID
	Sets     sql.NullInt32
	RepsMin  sql.NullInt32
	RepsMax  sql.NullInt32
	Weight   sql.NullFloat64
	Interval sql.NullInt32
	Rest     sql.NullInt32
}

func (q *Queries) EditWorkoutExercise(ctx context.Context, arg EditWorkoutExerciseParams) (WorkoutExercise, error) {
	row := q.db.QueryRowContext(ctx, editWorkoutExercise,
		arg.ID,
		arg.Sets,
		arg.RepsMin,
		arg.RepsMax,
		arg.Weight,
		arg.Interval,
		arg.Rest,
	)
	var i WorkoutExercise
	err := row.Scan(
		&i.ID,
		&i.WorkoutID,
		&i.ExerciseID,
		&i.Sets,
		&i.RepsMin,
		&i.RepsMax,
		&i.Weight,
		&i.Interval,
		&i.Rest,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
