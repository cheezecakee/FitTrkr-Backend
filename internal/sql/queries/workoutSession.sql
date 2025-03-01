-- name: GetActiveWorkoutSession :one
SELECT id
FROM workout_session
WHERE user_id = $1 AND is_active = true AND completed = false
LIMIT 1;

-- name: GetActiveExerciseSession :one
SELECT id, workout_exercise_id
FROM exercise_session
WHERE workout_session_id = $1 AND is_active = true AND completed = false
LIMIT 1;

-- name: CreateWorkoutSession :one
WITH existing_session AS (
    SELECT id
    FROM workout_session
    WHERE user_id = $1 AND is_active = true AND completed = false
    LIMIT 1
)
INSERT INTO workout_session (user_id, workout_id, is_active, expires_at)
SELECT $1, $2, true, NOW() + INTERVAL '6 hours'
WHERE NOT EXISTS (SELECT 1 FROM existing_session)
RETURNING id;

-- name: CreateExerciseSession :one
WITH active_session AS (
    SELECT id
    FROM workout_session
    WHERE user_id = $1 AND is_active = true AND completed = false
    LIMIT 1
)
INSERT INTO exercise_session (workout_session_id, workout_exercise_id, is_active)
SELECT (SELECT id FROM active_session), $2, true
WHERE EXISTS (SELECT 1 FROM active_session)
RETURNING id;

-- name: CompleteWorkoutSession :one
UPDATE workout_session
SET completed = true, is_active = false, ended_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING id;

-- name: CompleteExerciseSession :one
UPDATE exercise_session
SET completed = true, is_active = false, ended_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING id;

-- name: StopWorkoutSession :one
UPDATE workout_session
SET is_active = false, ended_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING id;

-- name: StopExerciseSession :one
UPDATE exercise_session
SET is_active = false, ended_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING id;

-- name: StopSession :one
BEGIN;

-- Stop the Workout Session
WITH active_sessions AS (
    SELECT es.id AS exercise_session_id
    FROM workout_session ws
    LEFT JOIN exercise_session es ON es.workout_session_id = ws.id
    WHERE ws.id = $1 AND ws.is_active = true
    AND es.is_active = true
    AND es.completed = false
)
-- Stop the Workout Session
UPDATE workout_session
SET is_active = false, ended_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING id;

-- If there's any active ExerciseSession, stop it as well
UPDATE exercise_session
SET is_active = false, ended_at = NOW()
WHERE workout_session_id = $1 AND is_active = true;

COMMIT;
