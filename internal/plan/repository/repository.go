package repository

import (
	"context"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/plan/models"
)

type PlanRepo interface {
	Create(ctx context.Context, plan *m.Plan) (*m.Plan, error)
	GetByID(ctx context.Context, id uint) (*m.Plan, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Plan, error)
	Update(ctx context.Context, plan *m.Plan) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*m.Plan, error)
}

type PlanExRepo interface {
	Create(ctx context.Context, planEx *m.PlanEx) (*m.PlanEx, error)
	GetByID(ctx context.Context, id uint) (*m.PlanEx, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*m.PlanEx, error)
	Update(ctx context.Context, planEx *m.PlanEx) error
	Delete(ctx context.Context, id uint) error
	DeleteByPlanID(ctx context.Context, planID uint) error
}
