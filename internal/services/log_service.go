package service

import (
	"context"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/log/model"
)

type LogService interface {
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

type DBLogService struct {
	repo repo.LogRepo
}

func NewLogService(repo repo.LogRepo) LogService {
	return &DBLogService{repo: repo}
}

func (s *DBLogService) Create(ctx context.Context, log *m.Log) error {
	return s.repo.Create(ctx, log)
}

func (s *DBLogService) GetByID(ctx context.Context, id uint) (*m.Log, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBLogService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Log, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *DBLogService) GetByPlanID(ctx context.Context, planID uint) ([]*m.Log, error) {
	return s.repo.GetByPlanID(ctx, planID)
}

func (s *DBLogService) GetByType(ctx context.Context, logType string) ([]*m.Log, error) {
	return s.repo.GetByType(ctx, logType)
}

func (s *DBLogService) GetByPriority(ctx context.Context, priority string) ([]*m.Log, error) {
	return s.repo.GetByPriority(ctx, priority)
}

func (s *DBLogService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *DBLogService) List(ctx context.Context, offset, limit int) ([]*m.Log, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *DBLogService) ListByRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Log, error) {
	return s.repo.ListByRange(ctx, userID, startDate, endDate)
}

func (s *DBLogService) GetPRs(ctx context.Context, userID uuid.UUID) ([]*m.Log, error) {
	return s.repo.GetPRs(ctx, userID)
}
