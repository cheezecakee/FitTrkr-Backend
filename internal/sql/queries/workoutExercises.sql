-- name: CreateWorkoutExercise :one
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
RETURNING *;

-- name: EditWorkoutExercise :one
UPDATE workout_exercises
SET sets = COALESCE(sqlc.narg('sets'), sets),
    reps_min = COALESCE(sqlc.narg('reps_min'), reps_min),
    reps_max = COALESCE(sqlc.narg('reps_max'), reps_max),
    weight = COALESCE(sqlc.narg('weight'), weight),
    interval = COALESCE(sqlc.narg('interval'), interval),
    rest = COALESCE(sqlc.narg('rest'), rest),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetWorkoutExercises :many
SELECT * FROM workout_exercises WHERE workout_id = $1;

-- name: DeleteWorkoutExercise :exec
DELETE FROM workout_exercises WHERE id = $1;
