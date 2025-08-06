package exercise

import (
	"context"
)

type EquipmentService interface {
	Create(ctx context.Context, equipment *Equipment) error
	GetByID(ctx context.Context, id int) (*Equipment, error)
	GetByName(ctx context.Context, name string) (*Equipment, error)
	List(ctx context.Context, offset, limit int) ([]*Equipment, error)
}

type DBEquipmentService struct {
	repo EquipmentRepo
}

func NewEquipmentService(repo EquipmentRepo) EquipmentService {
	return &DBEquipmentService{repo: repo}
}

func (s *DBEquipmentService) Create(ctx context.Context, equipment *Equipment) error {
	return s.repo.Create(ctx, equipment)
}

func (s *DBEquipmentService) GetByID(ctx context.Context, id int) (*Equipment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBEquipmentService) GetByName(ctx context.Context, name string) (*Equipment, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *DBEquipmentService) List(ctx context.Context, offset, limit int) ([]*Equipment, error) {
	return s.repo.List(ctx, offset, limit)
}
