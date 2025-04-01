package service

import (
	"context"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/repository"
)

type ExerciseService interface {
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

type DBExerciseService struct {
	repo repository.ExerciseRepo
}

func NewExerciseService(repo repository.ExerciseRepo) ExerciseService {
	return &DBExerciseService{repo: repo}
}

func (s *DBExerciseService) Create(ctx context.Context, exercise *m.Exercise) error {
	return s.repo.Create(ctx, exercise)
}

func (s *DBExerciseService) GetByID(ctx context.Context, id uint) (*m.Exercise, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBExerciseService) GetByName(ctx context.Context, name string) (*m.Exercise, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *DBExerciseService) GetByCategory(ctx context.Context, category string) ([]*m.Exercise, error) {
	return s.repo.GetByCategory(ctx, category)
}

func (s *DBExerciseService) GetByEquipment(ctx context.Context, equipment string) ([]*m.Exercise, error) {
	return s.repo.GetByEquipment(ctx, equipment)
}

func (s *DBExerciseService) Update(ctx context.Context, exercise *m.Exercise) error {
	return s.repo.Update(ctx, exercise)
}

func (s *DBExerciseService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *DBExerciseService) List(ctx context.Context, offset, limit int) ([]*m.Exercise, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *DBExerciseService) Search(ctx context.Context, query string) ([]*m.Exercise, error) {
	return s.repo.Search(ctx, query)
}
