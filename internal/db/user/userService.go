// Package user provides business logic and data access for users.
package user

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr/internal/utils/auth"
	"github.com/cheezecakee/fitrkr/internal/utils/helper"
)

var (
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrDuplicateUsername  = errors.New("username already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	Register(ctx context.Context, user User) (User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]User, error)
}

type userService struct {
	repo       UserRepo
	jwtManager auth.JWT
}

func NewUserService(repo UserRepo, jwtMgr auth.JWT) UserService {
	return &userService{repo: repo, jwtManager: jwtMgr}
}

func (s *userService) Register(ctx context.Context, user User) (User, error) {
	// Default roles to ["user"] if not provided
	if len(user.Roles) == 0 {
		user.Roles = []string{"user"}
	}
	// Check for existing Email
	u, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return User{}, err
	}
	if err == nil && u.ID != uuid.Nil {
		return User{}, ErrDuplicateEmail
	}

	// Check for existing Username
	u, err = s.repo.GetByUsername(ctx, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return User{}, err
	}
	if err == nil && u.ID != uuid.Nil {
		return User{}, ErrDuplicateUsername
	}

	hashedPassword, err := helper.HashPassword(user.PasswordHash)
	if err != nil {
		return User{}, err
	}
	user.PasswordHash = hashedPassword

	return s.repo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println("email err:", err)
		return "", ErrInvalidCredentials // No logging for security
	}
	log.Printf("email: %s, username: %s", user.Email, user.Username)

	if err := helper.ComparePassword(user.PasswordHash, password); err != nil {
		log.Println("password compare err:", err)
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtManager.MakeJWT(user.ID, user.Roles)
	if err != nil {
		log.Println("jwt err:", err)
		return "", err
	}

	return token, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, user User) (User, error) {
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return User{}, err
	}

	if user.PasswordHash != "" && user.PasswordHash != existingUser.PasswordHash {
		hashedPassword, err := helper.HashPassword(user.PasswordHash)
		if err != nil {
			return User{}, err
		}
		user.PasswordHash = hashedPassword
	} else {
		// Preserve existing hash if no new password is provided
		user.PasswordHash = existingUser.PasswordHash
	}
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, offset, limit int) ([]User, error) {
	if offset < 0 {
		offset = 0
	}
	limit = helper.Clamp(limit, 10, 100)
	users, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}
