package service

import (
	"context"

	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/repository"
)

type ExSetService interface {
	Create(ctx context.Context, set *m.ExSet) error
	CreateBatch(ctx context.Context, sets []*m.ExSet) error
	GetByID(ctx context.Context, id uint) (*m.ExSet, error)
	GetBySessionExID(ctx context.Context, sessionExID uuid.UUID) ([]*m.ExSet, error)
	Update(ctx context.Context, set *m.ExSet) error
	Delete(ctx context.Context, id uint) error
}

type DBExSetService struct {
	repo repository.ExSetRepo
}

func NewExsetService(repo repository.ExSetRepo) ExSetService {
	return &DBExSetService{repo: repo}
}

func (s *DBExSetService) Create(ctx context.Context, set *m.ExSet) error {
	return s.repo.Create(ctx, set)
}

func (s *DBExSetService) CreateBatch(ctx context.Context, sets []*m.ExSet) error {
	return s.repo.CreateBatch(ctx, sets)
}

func (s *DBExSetService) GetByID(ctx context.Context, id uint) (*m.ExSet, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBExSetService) GetBySessionExID(ctx context.Context, sessionExID uuid.UUID) ([]*m.ExSet, error) {
	return s.repo.GetBySessionExID(ctx, sessionExID)
}

func (s *DBExSetService) Update(ctx context.Context, set *m.ExSet) error {
	return s.repo.Update(ctx, set)
}

func (s *DBExSetService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
