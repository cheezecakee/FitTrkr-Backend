package service

import (
	"context"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/session/model"
)

type SessionExService interface {
	Create(ctx context.Context, sessionEx *m.SessionEx) (*m.SessionEx, error)
	GetByID(ctx context.Context, id uuid.UUID) (*m.SessionEx, error)
	GetBysessionID(ctx context.Context, sessionID uuid.UUID) ([]*m.SessionEx, error)
	Update(ctx context.Context, sessionEx *m.SessionEx) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBysessionID(ctx context.Context, sessionID uuid.UUID) error
}

type DBSessionExService struct {
	repo repo.SessionExRepo
}

func NewSessionExService(repo repo.SessionExRepo) SessionExService {
	return &DBSessionExService{repo: repo}
}

func (s *DBSessionExService) Create(ctx context.Context, sessionEx *m.SessionEx) (*m.SessionEx, error) {
	return s.repo.Create(ctx, sessionEx)
}

func (s *DBSessionExService) GetByID(ctx context.Context, id uuid.UUID) (*m.SessionEx, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBSessionExService) GetBysessionID(ctx context.Context, sessionID uuid.UUID) ([]*m.SessionEx, error) {
	return s.repo.GetBysessionID(ctx, sessionID)
}

func (s *DBSessionExService) Update(ctx context.Context, sessionEx *m.SessionEx) error {
	return s.repo.Update(ctx, sessionEx)
}

func (s *DBSessionExService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *DBSessionExService) DeleteBysessionID(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.DeleteBysessionID(ctx, sessionID)
}
