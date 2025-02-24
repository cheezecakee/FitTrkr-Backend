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

	// User routes
	mux.Handle("POST /api/user/register", http.HandlerFunc(apiCfg.RegisterUser))
	mux.Handle("POST /api/user/login", http.HandlerFunc(apiCfg.LoginUser))
	mux.Handle("POST /api/user/logout", http.HandlerFunc(apiCfg.LogoutUser))
	mux.Handle("POST /api/user/revoke", http.HandlerFunc(apiCfg.PostRevoke))
	mux.Handle("PUT /api/user/edit", http.HandlerFunc(apiCfg.EditUser))
	mux.Handle("DELETE /api/user", http.HandlerFunc(apiCfg.DeleteUser))
	mux.Handle("POST /api/user/refresh", http.HandlerFunc(apiCfg.PostRefresh))

	// Workout routes
	// mux.Handle("GET /api/workouts/{id}", http.HandlerFunc(GetWorkouts))
	mux.Handle("GET /api/workouts", http.HandlerFunc(apiCfg.GetWorkouts))
	mux.Handle("POST /api/workouts", http.HandlerFunc(apiCfg.CreateWorkout))
	mux.Handle("PUT /api/workouts/{id}", http.HandlerFunc(apiCfg.EditWorkout))
	mux.Handle("DELETE /api/workouts/{id}", http.HandlerFunc(apiCfg.DeleteWorkout))

	// Workout Exercises routes
	mux.Handle("POST /api/workouts/exercises", http.HandlerFunc(apiCfg.CreateWorkoutExercise))
	mux.Handle("PUT /api/workouts/exercises/{id}", http.HandlerFunc(apiCfg.EditWorkoutExercise))
	mux.Handle("DELETE /api/workouts/exercises/{id}", http.HandlerFunc(apiCfg.DeleteWorkoutExercise))

	// Exercise routes

	// Start/Stop Workout
	mux.Handle("POST /api/workouts/start", http.HandlerFunc(StartWorkout))
	mux.Handle("POST /api/workouts/stop", http.HandlerFunc(StopWorkout))

	// Logs
	mux.Handle("GET /api/logs", http.HandlerFunc(GetWorkoutLogs))

	return secureHeaders(mux)
}
