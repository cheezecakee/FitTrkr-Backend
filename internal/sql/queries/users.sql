-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, age, first_name, last_name)
VALUES (
    gen_random_uuid(), -- Generate a UUID
    NOW(), -- Set created_at to the current timestamp
    NOW(), -- Set updated_at to the current timestamp
    $1, -- email
    $2, -- age
    $3, -- first_name
    $4  -- last_name
)
RETURNING *;

-- name: RegisterUser :one
INSERT INTO users (id, email, password_hash, first_name, last_name, age, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,  -- email
    $2,  -- password_hash
    $3,  -- first_name
    $4,  -- last_name
    $5,  -- age
    NOW(), 
    NOW()
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: EditUser :one
UPDATE users
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    age = COALESCE($4, age),
    password_hash = COALESCE($5, password_hash),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

