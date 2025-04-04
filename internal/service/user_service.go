package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/repository"
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
	Register(ctx context.Context, user *m.User) (*m.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*m.User, error)
	GetUserByEmail(ctx context.Context, email string) (*m.User, error)
	GetUserByUsername(ctx context.Context, username string) (*m.User, error)
	Update(ctx context.Context, user *m.User) (*m.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*m.User, error)
}

type userService struct {
	repo       repository.UserRepo
	jwtManager auth.JWT
}

func NewUserService(repo repository.UserRepo, jwtMgr auth.JWT) UserService {
	return &userService{repo: repo, jwtManager: jwtMgr}
}

func (s *userService) Register(ctx context.Context, user *m.User) (*m.User, error) {
	// Check for existing Email
	if _, err := s.repo.GetByEmail(ctx, user.Email); err == nil {
		return nil, ErrDuplicateEmail
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	// Check for existing Username
	if _, err := s.repo.GetByUsername(ctx, user.Username); err == nil {
		return nil, ErrDuplicateUsername
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	hashedPassword, err := helper.HashPassword(user.PasswordHash)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashedPassword

	return s.repo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", ErrInvalidCredentials // No logging for security
		}
		return "", err
	}

	if err := helper.ComparePassword(user.PasswordHash, password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtManager.MakeJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*m.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*m.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*m.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, user *m.User) (*m.User, error) {
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if user.PasswordHash != existingUser.PasswordHash && existingUser.PasswordHash != "" {
		hashedPassword, err := helper.HashPassword(user.PasswordHash)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashedPassword
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

func (s *userService) List(ctx context.Context, offset, limit int) ([]*m.User, error) {
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
