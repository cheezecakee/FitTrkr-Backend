// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
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
RETURNING id, first_name, last_name, password_hash, email, age, created_at, updated_at, is_premium
`

type CreateUserParams struct {
	Email     string
	Age       sql.NullInt32
	FirstName string
	LastName  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.Age,
		arg.FirstName,
		arg.LastName,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PasswordHash,
		&i.Email,
		&i.Age,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsPremium,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const editUser = `-- name: EditUser :one
UPDATE users
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    age = COALESCE($4, age),
    password_hash = COALESCE($5, password_hash),
    updated_at = NOW()
WHERE id = $1
RETURNING id, first_name, last_name, password_hash, email, age, created_at, updated_at, is_premium
`

type EditUserParams struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Age          sql.NullInt32
	PasswordHash string
}

func (q *Queries) EditUser(ctx context.Context, arg EditUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, editUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Age,
		arg.PasswordHash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PasswordHash,
		&i.Email,
		&i.Age,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsPremium,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, first_name, last_name, password_hash, email, age, created_at, updated_at, is_premium FROM users WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PasswordHash,
		&i.Email,
		&i.Age,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsPremium,
	)
	return i, err
}

const registerUser = `-- name: RegisterUser :one
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
RETURNING id, first_name, last_name, password_hash, email, age, created_at, updated_at, is_premium
`

type RegisterUserParams struct {
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Age          sql.NullInt32
}

func (q *Queries) RegisterUser(ctx context.Context, arg RegisterUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, registerUser,
		arg.Email,
		arg.PasswordHash,
		arg.FirstName,
		arg.LastName,
		arg.Age,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.PasswordHash,
		&i.Email,
		&i.Age,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsPremium,
	)
	return i, err
}
