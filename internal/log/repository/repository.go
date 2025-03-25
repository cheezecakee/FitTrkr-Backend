package repository

import (
	"context"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/log/models"
)

type LogRepo interface {
	Create(ctx context.Context, log *m.Log) error
	GetByID(ctx context.Context, id uint) (*m.Log, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Log, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*m.Log, error)
	GetByType(ctx context.Context, logType string) ([]*m.Log, error)
	GetByPriority(ctx context.Context, priority string) ([]*m.Log, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*m.Log, error)
	ListByRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Log, error)
	GetPRs(ctx context.Context, userID uuid.UUID) ([]*m.Log, error)
}
