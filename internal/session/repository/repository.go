package repository

import (
	"context"

	"github.com/google/uuid"

	"github/cheezecakee/fitrkr/internal/session/models"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) (*models.Session, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error)
	GetByPlanID(ctx context.Context, plan uint) ([]*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.Session, error)
	ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*models.Session, error)
}

type SessionExerciseRepository interface {
	Create(ctx context.Context, sessionExercise *models.SessionExercise) (*models.SessionExercise, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.SessionExercise, error)
	GetBysessionID(ctx context.Context, sessionID uuid.UUID) ([]*models.SessionExercise, error)
	Update(ctx context.Context, sessionExercise *models.SessionExercise) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBysessionID(ctx context.Context, id uuid.UUID) error
}

type ExerciseSetRepository interface {
	Create(ctx context.Context, set *models.ExerciseSet) error
	CreateBatch(ctx context.Context, sets []*models.ExerciseSet) error
	GetByID(ctx context.Context, id uint) (*models.ExerciseSet, error)
	GetBySessionExerciseID(ctx context.Context, sessionExerciseID uuid.UUID) ([]*models.ExerciseSet, error)
	Update(ctx context.Context, set *models.ExerciseSet) error
	Delete(ctx context.Context, id uint) error
}
