package exercise

import (
	"context"
)

type TrainingTypeService interface {
	Create(ctx context.Context, exerciseType *TrainingType) error
	GetByID(ctx context.Context, id int) (*TrainingType, error)
	GetByName(ctx context.Context, name string) (*TrainingType, error)
	List(ctx context.Context, offset, limit int) ([]*TrainingType, error)
}

type trainingTypeService struct {
	repo TrainingTypeRepo
}

func NewTrainingTypeService(repo TrainingTypeRepo) TrainingTypeService {
	return &trainingTypeService{repo: repo}
}

func (s *trainingTypeService) Create(ctx context.Context, exerciseType *TrainingType) error {
	return s.repo.Create(ctx, exerciseType)
}

func (s *trainingTypeService) GetByID(ctx context.Context, id int) (*TrainingType, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *trainingTypeService) GetByName(ctx context.Context, name string) (*TrainingType, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *trainingTypeService) List(ctx context.Context, offset, limit int) ([]*TrainingType, error) {
	return s.repo.List(ctx, offset, limit)
}
