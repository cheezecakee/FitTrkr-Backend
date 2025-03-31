package router

import (
	"net/http"

	u "github.com/cheezecakee/go-backend-utils/pkg/util"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(version string) http.Handler {
	r := chi.NewRouter()

	r.Use(u.LoggingMiddleware)

	apiPath := "/api/v" + version

	versionedRoutes := map[string]http.Handler{
		"/users": SetupUserRoutes(),
	}

	// Mount the versioned routes correctly by prefixing the apiPath
	for path, handler := range versionedRoutes {
		r.Mount(apiPath+path, handler)
	}

	return r
}

func SetupUserRoutes() http.Handler {
	r := chi.NewRouter()

	// Public routes (No auth required)
	r.Post("/register")
	r.Post("/login")

	// Protected routes (Auth required)
	r.Group(func(r chi.Router) {
		r.Use() // Apply auth middleware to all routes inside this group
		r.Post("/logout")
		r.Put("/edit")
		r.Delete("/")
	})

	return r
}
