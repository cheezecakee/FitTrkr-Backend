package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/repository"
	h "github.com/cheezecakee/fitrkr/pkg/helper"
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
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &DBUserService{repo: repo}
}

func (s *DBUserService) Register(ctx context.Context, user *m.User) (*m.User, error) {
	// Check for existing Email
	existingUser, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, err
	}

	// Check for existing Username
	existingUser, err = s.repo.GetByUsername(ctx, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(ctx, user)
}

func (s *DBUserService) Login(ctx context.Context, email, password string) (*m.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err // No logging for security
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *DBUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*m.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

func (s *DBUserService) GetUserByEmail(ctx context.Context, email string) (*m.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

func (s *DBUserService) GetUserByUsername(ctx context.Context, username string) (*m.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

func (s *DBUserService) Update(ctx context.Context, user *m.User) (*m.User, error) {
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if user.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	} else {
		user.PasswordHash = existingUser.PasswordHash
	}

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *DBUserService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
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
		return nil, err
	}
	return users, nil
}
