package repository

import (
	"context"

	m "github/cheezecakee/fitrkr/internal/exercise/models"
)

type ExerciseRepo interface {
	Create(ctx context.Context, exercise *m.Exercise) error
	GetByID(ctx context.Context, id uint) (*m.Exercise, error)
	GetByName(ctx context.Context, name string) (*m.Exercise, error)
	GetByCategory(ctx context.Context, category string) ([]*m.Exercise, error)
	GetByEquipment(ctx context.Context, equipment string) ([]*m.Exercise, error)
	Update(ctx context.Context, exercise *m.Exercise) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*m.Exercise, error)
	Search(ctx context.Context, query string) ([]*m.Exercise, error)
}
