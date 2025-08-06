package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	_ "github.com/cheezecakee/fitrkr/internal/api/docs"
	"github.com/cheezecakee/fitrkr/internal/api/router"
	"github.com/cheezecakee/fitrkr/internal/app"
	"github.com/cheezecakee/fitrkr/internal/utils/config"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT-based authentication using Bearer token
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	cfg := config.LoadConfig()
	app := app.NewApp(cfg.DBConnString, cfg.JWTManager)
	defer app.DB.Close()

	mux := router.SetupRouter(app, cfg.JWTManager, "1") // Pass dbQueries to your router
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
