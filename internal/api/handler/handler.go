package handler

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/cheezecakee/fitrkr/internal/app"
	"github.com/cheezecakee/fitrkr/internal/service"
	"github.com/cheezecakee/fitrkr/internal/utils/auth"
	"github.com/cheezecakee/fitrkr/internal/utils/helper"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

type Api struct {
	UserSvc      service.UserService
	PlanSvc      service.PlanService
	PlanExSvc    service.PlanExService
	SessionSvc   service.SessionService
	SessionExSvc service.SessionExService
	ExSetSvc     service.ExSetService
	LogSvc       service.LogService
	ExerciseSvc  service.ExerciseService
	JWTManager   auth.JWT
	Helper       *helper.Helper
}

func NewApi(app *app.App, jwtMgr auth.JWT, helper *helper.Helper) *Api {
	return &Api{
		UserSvc:      app.UserSvc,
		PlanSvc:      app.PlanSvc,
		PlanExSvc:    app.PlanExSvc,
		SessionSvc:   app.SessionSvc,
		SessionExSvc: app.SessionExSvc,
		ExSetSvc:     app.ExSetSvc,
		LogSvc:       app.LogSvc,
		ExerciseSvc:  app.ExerciseSvc,
		JWTManager:   jwtMgr,
		Helper:       helper,
	}
}

// Errors
func (api *Api) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Println(trace)
	// logger.Log.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (api *Api) ClientError(w http.ResponseWriter, status int) {
	// logger.Log.InfoLog.Printf("Client error: %d", status)
	http.Error(w, http.StatusText(status), status)
}

func (api *Api) NotFound(w http.ResponseWriter) {
	api.ClientError(w, http.StatusNotFound)
}
