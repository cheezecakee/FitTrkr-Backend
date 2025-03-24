package repository

import (
	"context"

	"github/cheezecakee/fitrkr/internal/exercise/models"
)

type ExerciseRepository interface {
	Create(ctx context.Context, exercise *models.Exercise) error
	GetByID(ctx context.Context, id uint) (*models.Exercise, error)
	GetByName(ctx context.Context, name string) (*models.Exercise, error)
	GetByCategory(ctx context.Context, category string) ([]*models.Exercise, error)
	GetByEquipment(ctx context.Context, equipment string) ([]*models.Exercise, error)
	Update(ctx context.Context, exercise *models.Exercise) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Exercise, error)
	Search(ctx context.Context, query string) ([]*models.Exercise, error)
}
