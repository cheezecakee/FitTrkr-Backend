package service

import (
	"context"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/repository"
)

type PlanExService interface {
	Create(ctx context.Context, planEx *m.PlanEx) (*m.PlanEx, error)
	GetByID(ctx context.Context, id uint) (*m.PlanEx, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*m.PlanEx, error)
	Update(ctx context.Context, planEx *m.PlanEx) error
	Delete(ctx context.Context, id uint) error
	DeleteByPlanID(ctx context.Context, planID uint) error
}

type DBPlanExService struct {
	repo repository.PlanExRepo
}

func NewPlanExService(repo repository.PlanExRepo) PlanExService {
	return &DBPlanExService{repo: repo}
}

func (s *DBPlanExService) Create(ctx context.Context, planEx *m.PlanEx) (*m.PlanEx, error) {
	return s.repo.Create(ctx, planEx)
}

func (s *DBPlanExService) GetByID(ctx context.Context, id uint) (*m.PlanEx, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBPlanExService) GetByPlanID(ctx context.Context, planID uint) ([]*m.PlanEx, error) {
	return s.repo.GetByPlanID(ctx, planID)
}

func (s *DBPlanExService) Update(ctx context.Context, planEx *m.PlanEx) error {
	return s.repo.Update(ctx, planEx)
}

func (s *DBPlanExService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *DBPlanExService) DeleteByPlanID(ctx context.Context, planID uint) error {
	return s.repo.DeleteByPlanID(ctx, planID)
}
