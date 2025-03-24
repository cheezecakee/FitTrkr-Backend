package repository

import (
	"database/sql"

	exerciseRepo "github/cheezecakee/fitrkr/internal/exercise/repository"
	logRepo "github/cheezecakee/fitrkr/internal/log/repository"
	planRepo "github/cheezecakee/fitrkr/internal/plan/repository"
	sessionRepo "github/cheezecakee/fitrkr/internal/session/repository"
	userRepo "github/cheezecakee/fitrkr/internal/user/repository"
)

type RepositoryManager struct {
	UserRepo     userRepo.UserRepository
	PlanRepo     planRepo.PlanRepository
	SessionRepo  sessionRepo.SessionRepository
	LogRepo      logRepo.LogRepository
	ExerciseRepo exerciseRepo.ExerciseRepository
}

func NewRepositoryManager(db *sql.DB) *RepositoryManager {
	return &RepositoryManager{
		UserRepo:     userRepo.NewUserRepository(db),
		PlanRepo:     planRepo.NewPlanRepository(db),
		SessionRepo:  sessionRepo.NewSessionRepository(db),
		LogRepo:      logRepo.NewLogRepository(db),
		ExerciseRepo: exerciseRepo.NewExerciseRepository(db),
	}
}
