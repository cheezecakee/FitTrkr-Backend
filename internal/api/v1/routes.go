package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr/internal/api/handler"
)

func SetupRoutes(api *handler.Api) http.Handler {
	r := chi.NewRouter()

	versionedRoutes := map[string]http.Handler{
		"/users": SetupUserRoutes(api),
		"/auth":  SetupAuthRoutes(api),
		// "/plans":    SetupWorkoutRoutes(api),
		// "/sessions": SetupSessionRoutes(api),
		// "/logs":     SetupLogRoutes(api),
		// "/admin":    SetupAdminRoutes(api),
	}

	// Mount the versioned routes
	for path, handler := range versionedRoutes {
		r.Mount(path, handler)
	}

	return r
}

func SetupUserRoutes(api *handler.Api) http.Handler {
	r := chi.NewRouter()

	// Public routes (No auth required)
	r.Post("/", api.CreateUser)

	// Protected routes (Auth required)
	r.Group(func(r chi.Router) {
		r.Use(api.IsAuthenticated())  // Apply auth middleware to all routes inside this group
		r.Put("/", api.UpdateUser)    // Using jwt token
		r.Delete("/", api.DeleteUser) // Using jwt token
	})

	return r
}

func SetupAuthRoutes(api *handler.Api) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", api.Login)
	r.Post("/refresh", api.RefreshToken)
	r.Post("/revoke", api.RevokeToken)

	r.Group(func(r chi.Router) {
		r.Use(api.IsAuthenticated())
		r.Post("/logout", api.Logout)
	})

	return r
}
