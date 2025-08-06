package user

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/cheezecakee/fitrkr/internal/utils/transaction"
)

type UserRepo interface {
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]User, error)
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
    INSERT INTO users (username, first_name, last_name, password_hash, email, roles)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium, roles`

func (r *userRepo) Create(ctx context.Context, user User) (User, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	var newUser User
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, createUser, user.Username, user.FirstName, user.LastName, user.PasswordHash, user.Email, pq.Array(user.Roles)).Scan(&newUser.ID, &newUser.Username, &newUser.FirstName, &newUser.LastName, &newUser.PasswordHash, &newUser.Email, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.IsPremium, pq.Array(&newUser.Roles))
	})
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

const getUserByID = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium, roles FROM users WHERE id = $1 LIMIT 1`

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	user := User{}
	err := r.tx.DB().QueryRowContext(ctx, getUserByID, id).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
		pq.Array(&user.Roles),
	)
	if err != nil {
		log.Printf("GetByID failed for %s: %v", id, err) // debug
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}

	log.Printf("GetByID found: ID=%s, Email=%s", user.ID, user.Email) // Debug
	return user, nil
}

const getUserByEmail = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium, roles FROM users WHERE email = $1 LIMIT 1`

func (r *userRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	user := User{}
	err := r.tx.DB().QueryRowContext(ctx, getUserByEmail, email).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
		pq.Array(&user.Roles),
	)
	if err != nil {
		log.Printf("GetByEmail failed for %s: %v", email, err) // debug
		return User{}, err
	}
	log.Printf("GetByEmail found: ID=%s, Email=%s", user.ID, user.Email) // Debug
	return user, nil
}

const getUserByUsername = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium, roles FROM users WHERE username = $1 LIMIT 1`

func (r *userRepo) GetByUsername(ctx context.Context, username string) (User, error) {
	user := User{}
	err := r.tx.DB().QueryRowContext(ctx, getUserByUsername, username).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsPremium,
		pq.Array(&user.Roles),
	)
	if err != nil {
		log.Printf("GetByUsername failed for %s: %v", username, err) // debug
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}
	log.Printf("GetByUsername found: ID=%s, Username=%s", user.ID, user.Username) // Debug
	return user, nil
}

const updateUser = `
    UPDATE users 
    SET 
        first_name = COALESCE(NULLIF($2, ''), first_name),
        last_name = COALESCE(NULLIF($3, ''), last_name),
        password_hash = COALESCE(NULLIF($4, ''), password_hash),
        email = COALESCE(NULLIF($5, ''), email),
        updated_at = NOW(),
        is_premium = COALESCE($6, is_premium),
        roles = COALESCE($7, roles)
    WHERE id = $1
    RETURNING id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium, roles`

func (r *userRepo) Update(ctx context.Context, user User) (User, error) {
	var updatedUser User
	err := r.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, updateUser, user.ID, user.FirstName, user.LastName, user.PasswordHash, user.Email, user.IsPremium, pq.Array(user.Roles)).Scan(
			&updatedUser.ID,
			&updatedUser.Username,
			&updatedUser.FirstName,
			&updatedUser.LastName,
			&updatedUser.PasswordHash,
			&updatedUser.Email,
			&updatedUser.CreatedAt,
			&updatedUser.UpdatedAt,
			&updatedUser.IsPremium,
			pq.Array(&updatedUser.Roles),
		)
	})
	if err != nil {
		log.Printf("Update failed for ID %s: %v", user.ID, err)
		return User{}, err
	}
	return updatedUser, nil
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

const listUsers = `SELECT id, username, first_name, last_name, password_hash, email, created_at, updated_at, is_premium, roles FROM users OFFSET $1 LIMIT $2`

func (r *userRepo) List(ctx context.Context, offset, limit int) ([]User, error) {
	rows, err := r.tx.DB().QueryContext(ctx, listUsers, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
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
