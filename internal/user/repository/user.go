package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/user/models"
)

type DBUserRepo struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepo {
	return &DBUserRepo{
		db: db,
	}
}

const createUser = `INSERT INTO users (email, username, password_hash) Values ($1, $2, $3) returning id`

func (r *DBUserRepo) Create(ctx context.Context, user *m.User) (*m.User, error) {
	err := r.db.QueryRowContext(ctx, createUser, user.Email, user.Username, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

const getByID = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE id = $1 LIMIT 1`

func (r *DBUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*m.User, error) {
	user := &m.User{}
	row := r.db.QueryRowContext(ctx, getByID, id)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
	)

	return user, err
}

const getByEmail = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE email = $1 LIMIT 1`

func (r *DBUserRepo) GetByEmail(ctx context.Context, email string) (*m.User, error) {
	user := &m.User{}
	row := r.db.QueryRowContext(ctx, getByEmail, email)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
	)

	return user, err
}

const getByUsername = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE username = $1 LIMIT 1`

func (r *DBUserRepo) GetByUsername(ctx context.Context, username string) (*m.User, error) {
	user := &m.User{}
	row := r.db.QueryRowContext(ctx, getByUsername, username)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
	)

	return user, err
}

const updateUser = `UPDATE users 
SET 
    first_name = COALESCE(NULLIF($2, ''), first_name),
    last_name = COALESCE(NULLIF($3,  ''), last_name),
    username = COALESCE(NULLIF($4,  ''), username),
    password_hash = COALESCE(NULLIF($5, ''), password_hash),
    updated_at = NOW()
WHERE id = $1
RETURNING id, username, first_name, last_name, email, created_at, updated_at, is_premium`

func (r *DBUserRepo) Update(ctx context.Context, user *m.User) (*m.User, error) {
	row := r.db.QueryRowContext(ctx, updateUser, user.ID, user.Username, user.PasswordHash, user.FirstName, user.LastName, user.Email)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
	)

	return user, err
}

const deleteUser = `DELETE FROM users WHERE id = $1`

func (r *DBUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteUser, id)
	return err
}

const listUsers = `SELECT id, first_name, last_name, password_hash, email, age, created_at, updated_at, is_premium FROM users OFFSET $1 LIMIT $2`

func (r *DBUserRepo) List(ctx context.Context, offset, limit int) ([]*m.User, error) {
	rows, err := r.db.QueryContext(ctx, listUsers, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*m.User
	for rows.Next() {
		user := &m.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsPremium,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
