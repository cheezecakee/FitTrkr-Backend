package main

import (
	"log"
	"net/http"

	"github.com/cheezecakee/fitrkr/internal/api"
	"github.com/cheezecakee/fitrkr/internal/app"
	"github.com/cheezecakee/fitrkr/internal/utils/config"
)

func main() {
	cfg := config.LoadConfig()
	app := app.NewApp(cfg.DBConnString, cfg.Helper)
	defer app.DB.Close()

	mux := api.SetupRouter(app, cfg.JWTManager, cfg.Helper, "1") // Pass dbQueries to your router
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
