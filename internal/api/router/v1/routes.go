// Package v1 provides HTTP routing for version 1 of the API.
package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/cheezecakee/fitrkr/internal/api"
	"github.com/cheezecakee/fitrkr/internal/api/handler"
)

func SetupRoutes(api *api.API) http.Handler {
	r := chi.NewRouter()

	versionedRoutes := map[string]http.Handler{
		"/users":     SetupUserRoutes(api.UserH, api.AuthM),
		"/auth":      SetupAuthRoutes(api.AuthH, api.AuthM),
		"/playlists": SetupPlaylistRoutes(api.PlaylistH, api.AuthM),
		"/admin":     SetupAdminRoutes(api.ExerciseH, api.EquipmentH, api.ExerciseCategoryH, api.MuscleGroupH, api.TrainingTypeH, api.AuthM),
		"/swagger":   httpSwagger.WrapHandler,
	}

	// Mount the versioned routes
	for path, handler := range versionedRoutes {
		r.Mount(path, handler)
	}

	return r
}

func SetupUserRoutes(h *handler.UserHandler, authM *handler.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	// Public routes (No auth required)
	r.Post("/", h.CreateUser)

	// Protected routes (Auth required)
	r.Group(func(r chi.Router) {
		r.Use(authM.IsAuthenticated()) // Apply auth middleware to all routes inside this group
		r.Get("/me", h.GetCurrentUser) // GET /users/me - Get current user info
		r.Put("/", h.UpdateUser)       // PUT /users - Update user
		r.Delete("/", h.DeleteUser)    // DELETE /users - Delete user
	})

	return r
}

func SetupAuthRoutes(h *handler.AuthHandler, authM *handler.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/refresh", h.RefreshToken) // Currently not available
	r.Post("/revoke", h.RevokeToken)   // Currently not available

	r.Group(func(r chi.Router) {
		r.Use(authM.IsAuthenticated())
		r.Post("/logout", h.Logout)
	})

	return r
}

func SetupPlaylistRoutes(h *handler.PlaylistHandler, authM *handler.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	// All playlist routes require authentication
	r.Group(func(r chi.Router) {
		r.Use(authM.IsAuthenticated())

		// Main playlist CRUD operations
		r.Post("/", h.CreatePlaylist)       // POST /playlists
		r.Get("/", h.GetUserPlaylists)      // GET /playlists
		r.Get("/{id}", h.GetPlaylist)       // GET /playlists/{id}
		r.Put("/{id}", h.UpdatePlaylist)    // PUT /playlists/{id}
		r.Delete("/{id}", h.DeletePlaylist) // DELETE /playlists/{id}

		// Session-specific playlist data
		r.Get("/{id}/session", h.GetPlaylistForSession) // GET /playlists/{id}/session

		// Exercise management within playlists
		r.Post("/{id}/exercises", h.AddExerciseToPlaylist)        // POST /playlists/{id}/exercises
		r.Delete("/exercises/{id}", h.RemoveExerciseFromPlaylist) // DELETE /playlists/exercises/{id}

		// Block management within playlists
		r.Post("/{id}/blocks", h.CreateExerciseBlock) // POST /playlists/{id}/blocks

		// Reference data endpoints
		r.Get("/tags", h.GetTags) // GET /playlists/tags
	})

	return r
}

func SetupAdminRoutes(exerciseH *handler.ExerciseHandler, equipmentH *handler.EquipmentHandler, categoryH *handler.ExerciseCategoryHandler, muscleGroupH *handler.MuscleGroupHandler, exerciseTypeH *handler.TrainingTypeHandler, authM *handler.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(authM.IsAuthenticated())
		r.Route("/exercises", func(r chi.Router) {
			// Read-only routes for any authenticated user
			r.Get("/", exerciseH.List)
			r.Get("/{id}", exerciseH.GetByID)

			// Admin-only routes
			r.Group(func(r chi.Router) {
				r.Use(authM.RequireAdmin())
				r.Post("/", exerciseH.Create)
				r.Put("/{id}", exerciseH.Update)
				r.Delete("/{id}", exerciseH.Delete)
			})
		})

		r.Route("/equipment", func(r chi.Router) {
			// Read-only routes for any authenticated user
			r.Get("/", equipmentH.List)
			r.Get("/{id}", equipmentH.GetByID)
		})

		r.Route("/exercise-categories", func(r chi.Router) {
			// Read-only routes for any authenticated user
			r.Get("/", categoryH.List)
			r.Get("/{id}", categoryH.GetByID)
		})

		r.Route("/muscle-groups", func(r chi.Router) {
			r.Get("/", muscleGroupH.List)
		})

		r.Route("/exercise-types", func(r chi.Router) {
			r.Get("/", exerciseTypeH.List)
		})
	})
	return r
}
