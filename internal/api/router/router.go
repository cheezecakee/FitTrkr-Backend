package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr/internal/api"
	v1 "github.com/cheezecakee/fitrkr/internal/api/router/v1"
	"github.com/cheezecakee/fitrkr/internal/app"
	"github.com/cheezecakee/fitrkr/internal/utils/auth"
)

func SetupRouter(app *app.App, jwtMgr auth.JWT, version string) http.Handler {
	r := chi.NewRouter()

	r.Use()

	api := api.NewApi(app, jwtMgr)

	versions := []string{"1"}
	for _, version := range versions {
		apiPath := "/api/v" + version
		r.Mount(apiPath, v1.SetupRoutes(api))
	}

	return r
}
