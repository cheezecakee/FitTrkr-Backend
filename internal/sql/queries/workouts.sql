-- name: CreateWorkout :one
INSERT INTO workouts (id, user_id, name, description, created_at, updated_at)
VALUES (
    gen_random_uuid(), -- Generate a UUID
    $1, -- user_id (Foreign key from users table)
    $2, -- name (Workout name)
    $3, -- description (Optional)
    NOW(), -- Set created_at to the current timestamp
    NOW() -- Set updated_at to the current timestamp
)
RETURNING *;

-- name: EditWorkout :one
UPDATE workouts
SET name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: GetWorkout :many
SELECT * FROM workouts;

-- name: GetWorkoutsByID :many
SELECT * FROM workouts WHERE user_id = $1;

-- name: DeleteWorkout :exec
DELETE FROM workouts WHERE id = $1;

