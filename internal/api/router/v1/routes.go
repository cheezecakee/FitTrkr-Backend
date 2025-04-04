package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr/internal/api"
	"github.com/cheezecakee/fitrkr/internal/api/handler"
)

func SetupRoutes(api *api.Api) http.Handler {
	r := chi.NewRouter()

	versionedRoutes := map[string]http.Handler{
		"/users": SetupUserRoutes(api.UserH, api.AuthM),
		"/auth":  SetupAuthRoutes(api.AuthH, api.AuthM),
		// "/plans": SetupPlanRoutes(api.PlanH, api.PlanExH, api.AuthM),
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

func SetupUserRoutes(h *handler.UserHandler, authM *handler.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	// Public routes (No auth required)
	r.Post("/", h.CreateUser)

	// Protected routes (Auth required)
	r.Group(func(r chi.Router) {
		r.Use(authM.IsAuthenticated()) // Apply auth middleware to all routes inside this group
		r.Put("/", h.UpdateUser)       // Using jwt token
		r.Delete("/", h.DeleteUser)    // Using jwt token
	})

	return r
}

func SetupAuthRoutes(h *handler.AuthHandler, authM *handler.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/refresh", h.RefreshToken)
	r.Post("/revoke", h.RevokeToken)

	r.Group(func(r chi.Router) {
		r.Use(authM.IsAuthenticated())
		r.Post("/logout", h.Logout)
	})

	return r
}

func SetupPlanRoutes(h *handler.PlanHandler, hEx *handler.PlanExHandler, authM *handler.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(authM.IsAuthenticated())

	// Plan routes
	r.Get("/", h.GetPlans)
	r.Post("/", h.CreatePlan)

	r.Route("/{planID}", func(r chi.Router) {
		r.Get("/", h.GetPlan)
		r.Put("/", h.UpdatePlan)
		r.Delete("/", h.DeletePlan)

		// Exercise routes under a specific plan
		r.Get("/exercises", hEx.GetPlanExercises)
		r.Post("/exercises", hEx.CreatePlanExercise)

		r.Route("/exercises/{exerciseID}", func(r chi.Router) {
			r.Put("/", hEx.UpdatePlanExercise)
			r.Delete("/", hEx.DeletePlanExercise)
		})
	})

	return r
}
