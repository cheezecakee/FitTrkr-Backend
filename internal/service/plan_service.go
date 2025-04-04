package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/repository"
)

var (
	ErrInvalidPlan  = errors.New("invalid plan")
	ErrPlanNotFound = errors.New("plan not found")
)

type PlanService interface {
	Create(ctx context.Context, plan *m.Plan) (*m.Plan, error)
	GetByID(ctx context.Context, id uint) (*m.Plan, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Plan, error)
	Update(ctx context.Context, plan *m.Plan) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*m.Plan, error)
}

type DBPlanService struct {
	repo repository.PlanRepo
}

func NewPlanService(repo repository.PlanRepo) PlanService {
	return &DBPlanService{repo: repo}
}

func (s *DBPlanService) Create(ctx context.Context, plan *m.Plan) (*m.Plan, error) {
	return s.repo.Create(ctx, plan)
}

func (s *DBPlanService) GetByID(ctx context.Context, id uint) (*m.Plan, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBPlanService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Plan, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *DBPlanService) Update(ctx context.Context, plan *m.Plan) error {
	return s.repo.Update(ctx, plan)
}

func (s *DBPlanService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *DBPlanService) List(ctx context.Context, offset, limit int) ([]*m.Plan, error) {
	return s.repo.List(ctx, offset, limit)
}
