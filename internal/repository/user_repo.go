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

type userRepo struct {
	tx transaction.BaseRepository
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		tx: transaction.NewBaseRepository(db),
	}
}

const createUser = `
    INSERT INTO users (username, first_name, last_name, password_hash, email)
    Values ($1, $2, $3, $4, $5) 
    returning id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium`

func (r *userRepo) Create(ctx context.Context, user *m.User) (*m.User, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	var newUser m.User
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createUser, user.Username, user.FirstName, user.LastName, user.PasswordHash, user.Email).Scan(&newUser.ID, &newUser.Username, &newUser.FirstName, &newUser.LastName, &newUser.PasswordHash, &newUser.Email, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.IsPremium)
	})
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}

const getUserByID = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE id = $1 LIMIT 1`

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*m.User, error) {
	user := &m.User{}
	err := r.tx.DB().QueryRowContext(ctx, getUserByID, id).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
	)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return user, nil
}

const getUserByEmail = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE email = $1 LIMIT 1`

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*m.User, error) {
	user := &m.User{}
	err := r.tx.DB().QueryRowContext(ctx, getUserByEmail, email).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	return user, nil
}

const getUserByUsername = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users WHERE username = $1 LIMIT 1`

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*m.User, error) {
	user := &m.User{}
	err := r.tx.DB().QueryRowContext(ctx, getUserByUsername, username).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}
	return user, nil
}

const updateUser = `
    UPDATE users 
    SET 
        first_name = COALESCE(NULLIF($2, ''), first_name),
        last_name = COALESCE(NULLIF($3,  ''), last_name),
        username = COALESCE(NULLIF($4,  ''), username),
        password_hash = COALESCE(NULLIF($5, ''), password_hash),
        updated_at = NOW()
    WHERE id = $1
    RETURNING id, username, first_name, last_name, email, created_at, updated_at, is_premium
`

func (r *userRepo) Update(ctx context.Context, user *m.User) (*m.User, error) {
	var updatedUser m.User
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, updateUser,
			user.ID, user.FirstName, user.LastName, user.Username, user.PasswordHash, user.Email,
		).Scan(
			&updatedUser.ID, &updatedUser.Username, &updatedUser.FirstName, &updatedUser.LastName,
			&updatedUser.PasswordHash, &updatedUser.Email, &updatedUser.CreatedAt, &updatedUser.UpdatedAt,
			&updatedUser.IsPremium,
		)
	})
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

const deleteUser = `DELETE FROM users WHERE id = $1`

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, deleteUser, id)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

const listUsers = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium FROM users OFFSET $1 LIMIT $2`

func (r *userRepo) List(ctx context.Context, offset, limit int) ([]*m.User, error) {
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
			&user.PasswordHash,
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

	return users, rows.Err()
}
