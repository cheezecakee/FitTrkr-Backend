package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type UserRepo interface {
	Create(ctx context.Context, user *m.User) (*m.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*m.User, error)
	GetByEmail(ctx context.Context, email string) (*m.User, error)
	GetByUsername(ctx context.Context, username string) (*m.User, error)
	Update(ctx context.Context, user *m.User) (*m.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*m.User, error)
}

type DBUserRepo struct {
	tx transaction.BaseRepository
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &DBUserRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createUser = `INSERT INTO users (email, username, password_hash) Values ($1, $2, $3) returning id`

func (r *DBUserRepo) Create(ctx context.Context, user *m.User) (*m.User, error) {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createUser, user.Email, user.Username, user.PasswordHash).Scan(&user.ID)
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

const getUserByID = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE id = $1 LIMIT 1`

func (r *DBUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*m.User, error) {
	user := &m.User{}
	row := r.tx.DB().QueryRowContext(ctx, getUserByID, id)
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

const getUserByEmail = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE email = $1 LIMIT 1`

func (r *DBUserRepo) GetByEmail(ctx context.Context, email string) (*m.User, error) {
	user := &m.User{}
	row := r.tx.DB().QueryRowContext(ctx, getUserByEmail, email)
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

const getUserByUsername = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE username = $1 LIMIT 1`

func (r *DBUserRepo) GetByUsername(ctx context.Context, username string) (*m.User, error) {
	user := &m.User{}
	row := r.tx.DB().QueryRowContext(ctx, getUserByUsername, username)
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
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, updateUser, user.ID, user.Username, user.PasswordHash, user.FirstName, user.LastName, user.Email)
		return row.Scan(
			&user.ID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,

			&user.IsPremium,
		)
	})
	if err != nil {
		return nil, err
	}
	return user, err
}

const deleteUser = `DELETE FROM users WHERE id = $1`

func (r *DBUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteUser, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const listUsers = `SELECT id, first_name, last_name, password_hash, email, age, created_at, updated_at, is_premium FROM users OFFSET $1 LIMIT $2`

func (r *DBUserRepo) List(ctx context.Context, offset, limit int) ([]*m.User, error) {
	rows, err := r.tx.DB().QueryContext(ctx, listUsers, offset, limit)
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
