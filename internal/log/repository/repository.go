package repository

import (
	"context"

	"github.com/google/uuid"

	"github/cheezecakee/fitrkr/internal/log/models"
)

type LogRepository interface {
	Create(ctx context.Context, log *models.Log) error
	GetByID(ctx context.Context, id uint) (*models.Log, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Log, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*models.Log, error)
	GetByType(ctx context.Context, logType string) ([]*models.Log, error)
	GetByPriority(ctx context.Context, priority string) ([]*models.Log, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Log, error)
	ListByRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*models.Log, error)
	GetPRs(ctx context.Context, userID uuid.UUID) ([]*models.Log, error)
}
