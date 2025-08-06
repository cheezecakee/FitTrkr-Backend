// Package app contains the application setup and lifecycle logic for FitTrkr.
package app

import (
	"database/sql"

	"github.com/cheezecakee/fitrkr/internal/db"
	"github.com/cheezecakee/fitrkr/internal/db/exercise"
	"github.com/cheezecakee/fitrkr/internal/db/playlist"
	"github.com/cheezecakee/fitrkr/internal/db/user"
	"github.com/cheezecakee/fitrkr/internal/utils/auth"
)

type App struct {
	DB                  *sql.DB
	UserSvc             user.UserService
	ExerciseSvc         exercise.ExerciseService
	ExerciseCategorySvc exercise.CategoryService
	EquipmentSvc        exercise.EquipmentService
	MuscleGroupSvc      exercise.MuscleGroupService
	TrainingTypeSvc     exercise.TrainingTypeService

	// Playlist services
	PlaylistSvc playlist.PlaylistService
}

func NewApp(DBConnstring string, jwtMgr auth.JWT) *App {
	database := db.NewConnection(DBConnstring)

	// Exercise domain repositories
	userRepo := user.NewUserRepo(database)
	exerciseRepo := exercise.NewExerciseRepo(database)
	exerciseCategoryRepo := exercise.NewCategoryRepo(database)
	equipmentRepo := exercise.NewEquipmentRepo(database)
	muscleGroupRepo := exercise.NewMuscleGroupRepo(database)
	TrainingTypeRepo := exercise.NewTrainingTypeRepo(database)

	// Playlist domain repositories
	playlistRepo := playlist.NewPlaylistRepo(database)
	exerciseBlockRepo := playlist.NewBlockRepo(database)
	playlistExerciseRepo := playlist.NewPlaylistExerciseRepo(database)
	exerciseConfigRepo := playlist.NewConfigRepo(database)

	// Initialize services
	playlistSvc := playlist.NewPlaylistService(
		playlistRepo,
		exerciseBlockRepo,
		playlistExerciseRepo,
		exerciseConfigRepo,
	)

	return &App{
		DB:                  database,
		UserSvc:             user.NewUserService(userRepo, jwtMgr),
		ExerciseSvc:         exercise.NewExerciseService(exerciseRepo),
		ExerciseCategorySvc: exercise.NewCategoryService(exerciseCategoryRepo),
		EquipmentSvc:        exercise.NewEquipmentService(equipmentRepo),
		MuscleGroupSvc:      exercise.NewMuscleGroupService(muscleGroupRepo),
		TrainingTypeSvc:     exercise.NewTrainingTypeService(TrainingTypeRepo),

		// Playlist service
		PlaylistSvc: playlistSvc,
	}
}
