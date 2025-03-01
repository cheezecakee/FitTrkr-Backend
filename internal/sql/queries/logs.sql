-- name: CreateSessionLog :one
INSERT INTO session_logs (user_id, workout_session_id, workout_exercise_logs_id, log_type, log_priority, log_message)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, workout_session_id, workout_exercise_logs_id, log_type, log_priority, log_message, created_at;

-- name: GetExerciseLogsForSession :many
SELECT sel.id, sel.set_number, sel.reps, sel.weight, sel.interval_sec, sel.notes, sel.created_at
FROM workout_exercise_logs sel
JOIN exercise_session wes ON wes.id = sel.exercise_session_id
WHERE wes.workout_session_id = $1
ORDER BY sel.set_number;

-- name: CreateExerciseLog :one
INSERT INTO workout_exercise_logs (exercise_session_id, set_number, reps, weight, interval_sec, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, exercise_session_id, set_number, reps, weight, interval_sec, notes, created_at;
