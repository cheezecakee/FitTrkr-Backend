package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	m "github/cheezecakee/fitrkr/internal/user/model"
	"github/cheezecakee/fitrkr/pkg/errors"
	h "github/cheezecakee/fitrkr/pkg/helper"
)

type UserService interface {
	Register(ctx context.Context, user *m.User) (*m.User, error)
	Login(ctx context.Context, email, password string) (*m.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*m.User, error)
	GetUserByEmail(ctx context.Context, email string) (*m.User, error)
	GetUserByUsername(ctx context.Context, username string) (*m.User, error)
	Update(ctx context.Context, user *m.User) (*m.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*m.User, error)
}

type DBUserService struct {
	repo repo.UserRepo
}

func NewUserService(repo repo.UserRepo) UserService {
	return &DBUserService{repo: repo}
}

func (s *DBUserService) Register(ctx context.Context, user *m.User) (*m.User, error) {
	// Check for existing Email
	existingUser, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.ErrInternalServer
	}
	if existingUser != nil {
		return nil, errors.ErrEmailExists
	}

	// Check for existing Username
	existingUser, err = s.repo.GetByUsername(ctx, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.ErrInternalServer
	}
	if existingUser != nil {
		return nil, errors.ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(ctx, user)
}

func (s *DBUserService) Login(ctx context.Context, email, password string) (*m.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.ErrInvalidCredentials // No logging for security
		}
		return nil, errors.ErrInternalServer
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	return user, nil
}

func (s *DBUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*m.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return user, nil
}

func (s *DBUserService) GetUserByEmail(ctx context.Context, email string) (*m.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return user, nil
}

func (s *DBUserService) GetUserByUsername(ctx context.Context, username string) (*m.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return user, nil
}

func (s *DBUserService) Update(ctx context.Context, user *m.User) (*m.User, error) {
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrInternalServer
	}

	if user.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		user.PasswordHash = string(hashedPassword)
	} else {
		user.PasswordHash = existingUser.PasswordHash
	}

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return updatedUser, nil
}

func (s *DBUserService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		return errors.ErrInternalServer
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (s *DBUserService) List(ctx context.Context, offset, limit int) ([]*m.User, error) {
	if offset < 0 {
		offset = 0
	}
	limit = h.Clamp(limit, 10, 100)
	users, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return users, nil
}
