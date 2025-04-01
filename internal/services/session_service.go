package service

import (
	"context"

	"github.com/google/uuid"

	m "github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/repository"
)

type SessionService interface {
	Create(ctx context.Context, session *m.Session) (*m.Session, error)
	GetByID(ctx context.Context, id uuid.UUID) (*m.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Session, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*m.Session, error)
	Update(ctx context.Context, session *m.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*m.Session, error)
	ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Session, error)
}

type DBSessionService struct {
	repo repository.SessionRepo
}

func NewSessionService(repo repository.SessionRepo) SessionService {
	return &DBSessionService{repo: repo}
}

func (s *DBSessionService) Create(ctx context.Context, session *m.Session) (*m.Session, error) {
	return s.repo.Create(ctx, session)
}

func (s *DBSessionService) GetByID(ctx context.Context, id uuid.UUID) (*m.Session, error) {
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *DBSessionService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Session, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *DBSessionService) GetByPlanID(ctx context.Context, planID uint) ([]*m.Session, error) {
	return s.repo.GetByPlanID(ctx, planID)
}

func (s *DBSessionService) Update(ctx context.Context, session *m.Session) error {
	return s.repo.Update(ctx, session)
}

func (s *DBSessionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *DBSessionService) List(ctx context.Context, offset, limit int) ([]*m.Session, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *DBSessionService) ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Session, error) {
	return s.repo.ListByDateRange(ctx, userID, startDate, endDate)
}
