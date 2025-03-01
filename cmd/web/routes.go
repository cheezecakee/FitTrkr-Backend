package main

import (
	"net/http"

	"github.com/cheezecakee/FitLogr/internal/database"
)

func (apiCfg *ApiConfig) SetupRoutes(dbQueries *database.Queries) http.Handler {
	mux := http.NewServeMux()

	// Middleware
	// mux.Use(middleware.LoggingMiddleware)

	mux.Handle("/internal/", http.StripPrefix("/internal", http.FileServer(http.Dir("./internal/"))))

	// Admin routes
	mux.Handle("GET /api/admin/user", http.HandlerFunc(apiCfg.GetUsers))

	// User routes
	mux.Handle("POST /api/user/register", http.HandlerFunc(apiCfg.RegisterUser))
	mux.Handle("POST /api/user/login", http.HandlerFunc(apiCfg.LoginUser))
	mux.Handle("POST /api/user/logout", apiCfg.ValidateSession(apiCfg.LogoutUser))
	mux.Handle("PUT /api/user/edit", apiCfg.ValidateSession(apiCfg.EditUser))
	mux.Handle("POST /api/user/revoke", http.HandlerFunc(apiCfg.PostRevoke))
	mux.Handle("POST /api/user/refresh", apiCfg.isAuthenticated(apiCfg.PostRefresh))
	mux.Handle("DELETE /api/user", apiCfg.ValidateSession(apiCfg.DeleteUser))

	// Workout routes
	mux.Handle("GET /api/workouts", apiCfg.ValidateSession(apiCfg.GetWorkouts))
	mux.Handle("POST /api/workouts", apiCfg.ValidateSession(apiCfg.CreateWorkout))
	mux.Handle("PUT /api/workouts/{id}", apiCfg.ValidateSession(apiCfg.EditWorkout))
	mux.Handle("DELETE /api/workouts/{id}", apiCfg.ValidateSession(apiCfg.DeleteWorkout))

	// Workout Exercises routes
	mux.Handle("POST /api/workouts/exercises", http.HandlerFunc(apiCfg.CreateWorkoutExercise))
	mux.Handle("PUT /api/workouts/exercises/{id}", http.HandlerFunc(apiCfg.EditWorkoutExercise))
	mux.Handle("DELETE /api/workouts/exercises/{id}", http.HandlerFunc(apiCfg.DeleteWorkoutExercise))

	// Start/Stop Workout
	mux.Handle("GET /api/session", apiCfg.ValidateSession(apiCfg.GetSession)) // Get current session details

	mux.Handle("GET /api/session/workout", apiCfg.ValidateSession(apiCfg.GetSessionWorkout)) // Get details of the ongoing workout
	mux.Handle("POST /api/session/workout/{id}", apiCfg.ValidateSession(apiCfg.PostSessionWorkout))

	mux.Handle("POST /api/session/workout/exercise/{id}/start", apiCfg.ValidateSession(apiCfg.StartExercise)) // Start first set (creates session if needed)

	mux.Handle("POST /api/session/start", apiCfg.ValidateSession(apiCfg.StopActiveWorkoutSession))   // User presses "Start"
	mux.Handle("PUT /api/session/stop", apiCfg.ValidateSession(apiCfg.StopActiveWorkoutSession))     // Stop full workout session early
	mux.Handle("PUT /api/session/finish", apiCfg.ValidateSession(apiCfg.FinishActiveWorkoutSession)) // Finish workout session

	// Start/Stop Workout
	mux.Handle("PUT /api/session/workout/finish", apiCfg.ValidateSession(apiCfg.FinishWorkoutSession)) // Mark workout completed

	mux.Handle("PUT /api/session/workout/exercise/log", apiCfg.ValidateSession(apiCfg.LogExerciseSet))    // Log set data
	mux.Handle("PUT /api/session/workout/exercise/stop", apiCfg.ValidateSession(apiCfg.StopExercise))     // Stop this exercise early
	mux.Handle("PUT /api/session/workout/exercise/finish", apiCfg.ValidateSession(apiCfg.FinishExercise)) // Mark exercise completed

	// Logs
	mux.Handle("GET /api/logs", http.HandlerFunc(GetWorkoutLogs))

	return secureHeaders(mux)
}
