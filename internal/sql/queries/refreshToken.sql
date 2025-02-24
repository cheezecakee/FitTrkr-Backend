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

-- name: GetUserFromRefreshToken :one
SELECT user_id, token, expires_at, revoked_at, created_at, updated_at FROM refresh_tokens WHERE token = $1 LIMIT 1;

-- name: RevokeRefreshTokenFromUser :exec
UPDATE refresh_tokens SET revoked_at = now(), updated_at = now() WHERE token = $1;
