package app

import (
	"database/sql"

	"github.com/cheezecakee/fitrkr/internal/repository"
	"github.com/cheezecakee/fitrkr/internal/service"
	"github.com/cheezecakee/fitrkr/internal/utils/helper"
)

type App struct {
	DB           *sql.DB
	UserSvc      service.UserService
	PlanSvc      service.PlanService
	PlanExSvc    service.PlanExService
	SessionSvc   service.SessionService
	SessionExSvc service.SessionExService
	ExSetSvc     service.ExSetService
	LogSvc       service.LogService
	ExerciseSvc  service.ExerciseService
}

func NewApp(DBConnstring string, helper *helper.Helper) *App {
	db := repository.NewDB(DBConnstring)

	// Initialize repos
	UserRepo := repository.NewUserRepo(db)
	PlanRepo := repository.NewPlanRepo(db)
	PlanExRepo := repository.NewPlanExRepo(db)
	SessionRepo := repository.NewSessionRepo(db)
	SessionExRepo := repository.NewSessionExRepo(db)
	ExSetRepo := repository.NewExSetRepo(db)
	LogRepo := repository.NewLogRepo(db)
	ExerciseRepo := repository.NewExerciseRepo(db)

	return &App{
		DB:           db,
		UserSvc:      service.NewUserService(UserRepo, helper),
		PlanSvc:      service.NewPlanService(PlanRepo),
		PlanExSvc:    service.NewPlanExService(PlanExRepo),
		SessionSvc:   service.NewSessionService(SessionRepo),
		SessionExSvc: service.NewSessionExService(SessionExRepo),
		ExSetSvc:     service.NewExsetService(ExSetRepo),
		LogSvc:       service.NewLogService(LogRepo),
		ExerciseSvc:  service.NewExerciseService(ExerciseRepo),
	}
}
