package repository

import (
	"context"

	"github.com/google/uuid"

	"github/cheezecakee/fitrkr/internal/user/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, user *models.User) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
}
