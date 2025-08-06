package exercise

import (
	"context"
)

type MuscleGroupService interface {
	Create(ctx context.Context, muscleGroup *MuscleGroup) error
	GetByID(ctx context.Context, id int) (*MuscleGroup, error)
	GetByName(ctx context.Context, name string) (*MuscleGroup, error)
	List(ctx context.Context, offset, limit int) ([]*MuscleGroup, error)
}

type muscleGroupService struct {
	repo MuscleGroupRepo
}

func NewMuscleGroupService(repo MuscleGroupRepo) MuscleGroupService {
	return &muscleGroupService{repo: repo}
}

func (s *muscleGroupService) Create(ctx context.Context, muscleGroup *MuscleGroup) error {
	return s.repo.Create(ctx, muscleGroup)
}

func (s *muscleGroupService) GetByID(ctx context.Context, id int) (*MuscleGroup, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *muscleGroupService) GetByName(ctx context.Context, name string) (*MuscleGroup, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *muscleGroupService) List(ctx context.Context, offset, limit int) ([]*MuscleGroup, error) {
	return s.repo.List(ctx, offset, limit)
}
