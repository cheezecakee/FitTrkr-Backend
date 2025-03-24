package repository

import (
	"context"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/user/models"
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
