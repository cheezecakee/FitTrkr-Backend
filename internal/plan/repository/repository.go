package repository

import (
	"context"

	"github.com/google/uuid"

	"github/cheezecakee/fitrkr/internal/Plan/models"
)

type PlanRepository interface {
	Create(ctx context.Context, plan *models.Plan) error
	GetByID(ctx context.Context, id uint) (*models.Plan, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Plan, error)
	Update(ctx context.Context, plan *models.Plan) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Plan, error)
}

type PlanExerciseRepository interface {
	Create(ctx context.Context, planExercise *models.PlanExercise) error
	GetByID(ctx context.Context, id uint) (*models.PlanExercise, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*models.PlanExercise, error)
	Update(ctx context.Context, planExercise *models.PlanExercise) error
	Delete(ctx context.Context, id uint) error
	DeleteByPlanID(ctx context.Context, planID uint) error
}
