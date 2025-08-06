package api

import (
	"github.com/cheezecakee/fitrkr/internal/api/handler"
	"github.com/cheezecakee/fitrkr/internal/app"
	"github.com/cheezecakee/fitrkr/internal/utils/auth"
)

type API struct {
	AuthH             *handler.AuthHandler
	AuthM             *handler.AuthMiddleware
	EquipmentH        *handler.EquipmentHandler
	ExerciseCategoryH *handler.ExerciseCategoryHandler
	ExerciseH         *handler.ExerciseHandler
	TrainingTypeH     *handler.TrainingTypeHandler
	MuscleGroupH      *handler.MuscleGroupHandler
	PlaylistH         *handler.PlaylistHandler
	UserH             *handler.UserHandler
}

func NewAPI(app *app.App, jwtMgr auth.JWT) *API {
	return &API{
		AuthH:             handler.NewAuthHandler(app.UserSvc),
		AuthM:             handler.NewAuthMiddleware(jwtMgr, app.UserSvc),
		EquipmentH:        handler.NewEquipmentHandler(app.EquipmentSvc),
		ExerciseCategoryH: handler.NewExerciseCategoryHandler(app.ExerciseCategorySvc),
		ExerciseH:         handler.NewExerciseHandler(app.ExerciseSvc),
		TrainingTypeH:     handler.NewTrainingTypeHandler(app.TrainingTypeSvc),
		MuscleGroupH:      handler.NewMuscleGroupHandler(app.MuscleGroupSvc),
		PlaylistH:         handler.NewPlaylistHandler(app.PlaylistSvc),
		UserH:             handler.NewUserHandler(app.UserSvc),
	}
}
