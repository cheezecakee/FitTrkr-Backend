-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at, created_at, updated_at, revoked_at)
VALUES (
    $1,
    $2,
    $3,
    now(),
    now(),
    NULL
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM refresh_tokens WHERE token = $1 LIMIT 1;

-- name: GetLatestSessionByID :one
SELECT * FROM refresh_tokens WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET is_revoked = true, revoked_at = now(), updated_at = now() WHERE token = $1;

-- name: DeleteSession :exec
DELETE FROM refresh_tokens WHERE user_id = $1;
