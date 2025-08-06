// Package router provides HTTP routing for the API.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr/internal/api"
	"github.com/cheezecakee/fitrkr/internal/api/handler"
	v1 "github.com/cheezecakee/fitrkr/internal/api/router/v1"
	"github.com/cheezecakee/fitrkr/internal/app"
	"github.com/cheezecakee/fitrkr/internal/utils/auth"
)

func SetupRouter(app *app.App, jwtMgr auth.JWT, version string) http.Handler {
	r := chi.NewRouter()

	r.Use(handler.CORS)

	api := api.NewAPI(app, jwtMgr)

	versions := []string{"1"}
	for _, version := range versions {
		apiPath := "/api/v" + version
		r.Mount(apiPath, v1.SetupRoutes(api))
	}

	return r
}
