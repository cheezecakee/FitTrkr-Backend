package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger/v2"

	_ "github.com/cheezecakee/FitLogr/internal/docs"
	"github.com/cheezecakee/FitLogr/internal/handlers"
)

func SetupRoutes(cfg *handler.Config, version string) http.Handler {
	r := chi.NewRouter()

	r.Use(cfg.LoggingMiddleware())

	apiPath := "/api/v" + version
	cfg.APIRoute = apiPath + "/"

	// Serve the swagger.json file from ./internal/docs
	fs := http.FileServer(http.Dir("./internal/docs"))
	r.Get(apiPath+"/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(apiPath+"/swagger/", fs).ServeHTTP(w, r)
	})

	// Set up Swagger UI, pointing it to the correct swagger.json URL

	r.Get(apiPath+"/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(apiPath+"swagger/swagger.json"),
	))

	versionedRoutes := map[string]http.Handler{
		"/users":    SetupUserRoutes(cfg),
		"/workouts": SetupWorkoutRoutes(cfg),
		"/sessions": SetupSessionRoutes(cfg),
		"/logs":     SetupLogRoutes(cfg),
		"/admin":    SetupAdminRoutes(cfg),
		"/auth":     SetupAuthRoutes(cfg),
	}

	// Mount the versioned routes correctly by prefixing the apiPath
	for path, handler := range versionedRoutes {
		r.Mount(apiPath+path, handler)
	}

	return r
}

func SetupUserRoutes(cfg *handler.Config) http.Handler {
	r := chi.NewRouter()

	// Public routes (No auth required)
	r.Post("/register", cfg.RegisterUser)
	r.Post("/login", cfg.LoginUser)

	// Protected routes (Auth required)
	r.Group(func(r chi.Router) {
		r.Use(cfg.IsAuthenticated()) // Apply auth middleware to all routes inside this group
		r.Post("/logout", cfg.LogoutUser)
		r.Put("/edit", cfg.EditUser)
		r.Delete("/", cfg.DeleteUser)
	})

	return r
}

func SetupWorkoutRoutes(cfg *handler.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(cfg.IsAuthenticated())

	r.Route("/", func(r chi.Router) {
		r.Get("/", cfg.GetWorkouts)
		r.Post("/", cfg.CreateWorkout)
	})

	r.Route("/{id}", func(r chi.Router) { // Workout ID
		r.Put("/", cfg.EditWorkout)
		r.Delete("/", cfg.DeleteWorkout)

		// Fetch workout exercises
		r.Get("/exercises", cfg.GetWorkoutsExercises)
		// Create workout exercise
		r.Post("/exercises", cfg.CreateWorkoutExercise)
	})

	r.Route("/exercise/{id}", func(r chi.Router) { // Exercise ID
		r.Put("/", cfg.EditWorkoutExercise)
		r.Delete("/", cfg.DeleteWorkoutExercise)
	})

	return r
}

func SetupSessionRoutes(cfg *handler.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(cfg.IsAuthenticated())

	// High-level session routes
	r.Get("/", cfg.GetSession)
	r.Get("/workouts", cfg.ListSessionWorkouts)
	r.Post("/workouts/{id}/start", cfg.StartSessionWorkout)
	r.Get("/workouts/{id}/exercises", cfg.ViewSessionWorkoutExercises) // Single tap on a workout to view contents

	// Grouped routes under /sessions/{sessionID}
	r.Route("/{sessionID}", func(r chi.Router) {
		r.Get("/exercises", cfg.ListSessionExercises)
		r.Post("/exercises/{id}/start", cfg.StartSessionExercise)
		r.Get("/exercise/{id}", cfg.ViewSessionExerciseDetails) // Single tap on an exercise to view contents
		r.Route("/{exerciseSessionID}", func(r chi.Router) {
			r.Post("/exercises", cfg.LogSet) // Log set
		})
	})

	return r
}

func SetupLogRoutes(cfg *handler.Config) http.Handler {
	r := chi.NewRouter()

	r.With(cfg.IsAuthenticated()).Get("/", cfg.GetLogs)

	return r
}

func SetupExercisesRoutes(cfg *handler.Config) http.Handler {
	return chi.NewRouter() // Placeholder for future endpoints
}

func SetupAdminRoutes(cfg *handler.Config) http.Handler {
	r := chi.NewRouter()

	r.Get("/users", cfg.GetUsers)
	r.Get("/version", cfg.GetVersion)

	return r
}

func SetupAuthRoutes(cfg *handler.Config) http.Handler {
	r := chi.NewRouter()

	r.Post("/revoke", cfg.PostRevoke)

	return r
}
