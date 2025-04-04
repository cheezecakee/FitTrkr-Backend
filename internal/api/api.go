package api

import (
	"github.com/cheezecakee/fitrkr/internal/api/handler"
	"github.com/cheezecakee/fitrkr/internal/app"
	"github.com/cheezecakee/fitrkr/internal/utils/auth"
)

type Api struct {
	UserH *handler.UserHandler
	AuthH *handler.AuthHandler
	// PlanH *handler.PlanHandler
	// PlanExH    *handler.PlanExHandler
	// SessionH   *handler.SessionHandler
	// SessionExH *handler.SessionExHandler
	// ExSetH     *handler.ExSetHandler
	// LogH       *handler.LogHandler
	// ExerciseH  *handler.ExerciseHandler
	AuthM *handler.AuthMiddleware
}

func NewApi(app *app.App, jwtMgr auth.JWT) *Api {
	return &Api{
		UserH: handler.NewUserHandler(app.UserSvc),
		AuthH: handler.NewAuthHandler(app.UserSvc),
		// PlanH: handler.NewPlanHandler(app.PlanSvc),
		// PlanExH:    handler.NewPlanExHandler(app.PlanExSvc),
		// SessionH:   handler.NewSessionHandler(app.SessionSvc),
		// SessionExH: handler.NewSessionExHandler(app.SessionExSvc),
		// ExSetH:     handler.NewExSetHandler(app.ExSetSvc),
		// LogH:       handler.NewLogHandler(app.LogSvc),
		// ExerciseH:  handler.NewExerciseHandler(app.ExerciseSvc),
		AuthM: handler.NewAuthMiddleware(jwtMgr),
	}
}
