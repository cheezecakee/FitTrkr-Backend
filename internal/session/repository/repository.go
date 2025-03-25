package repository

import (
	"context"

	"github.com/google/uuid"

	m "github/cheezecakee/fitrkr/internal/session/models"
)

type SessionRepo interface {
	Create(ctx context.Context, session *m.Session) (*m.Session, error)
	GetByID(ctx context.Context, id uuid.UUID) (*m.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*m.Session, error)
	GetByPlanID(ctx context.Context, planID uint) ([]*m.Session, error)
	Update(ctx context.Context, session *m.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*m.Session, error)
	ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*m.Session, error)
}

type SessionExRepo interface {
	Create(ctx context.Context, sessionEx *m.SessionEx) (*m.SessionEx, error)
	GetByID(ctx context.Context, id uuid.UUID) (*m.SessionEx, error)
	GetBysessionID(ctx context.Context, sessionID uuid.UUID) ([]*m.SessionEx, error)
	Update(ctx context.Context, sessionEx *m.SessionEx) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBysessionID(ctx context.Context, id uuid.UUID) error
}

type ExSetRepo interface {
	Create(ctx context.Context, set *m.ExSet) error
	CreateBatch(ctx context.Context, sets []*m.ExSet) error
	GetByID(ctx context.Context, id uint) (*m.ExSet, error)
	GetBySessionExID(ctx context.Context, sessionExID uuid.UUID) ([]*m.ExSet, error)
	Update(ctx context.Context, set *m.ExSet) error
	Delete(ctx context.Context, id uint) error
}
