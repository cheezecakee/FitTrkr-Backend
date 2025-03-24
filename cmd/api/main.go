package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	_ "github.com/cheezecakee/FitLogr/internal/docs"
	"github.com/cheezecakee/FitLogr/internal/handlers"
	r "github.com/cheezecakee/FitLogr/internal/router"
	"github.com/cheezecakee/FitLogr/internal/services"
)

// @title FitLogr API
// @version 1.0
// @description API for FitLogr
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer db.Close() // Ensure the database connection is closed when the program exits

	// Create database queries instance
	dbQueries := services.New(db)

	// Initialize jwt secret
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if dbURL == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	// Get Redis address from .env
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default Redis address
	}

	cfg := handler.NewConfig(dbQueries, jwtSecret, redisAddr)
	// Setup routes
	router := r.SetupRoutes(cfg, "1") // Pass dbQueries to your router

	// Start server
	port := ":8080"

	srv := &http.Server{
		Addr:     port,
		ErrorLog: cfg.Logger.ErrorLog,
		Handler:  router,
	}

	cfg.Logger.InfoLog.Printf("Server is running on https://localhost" + port)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	cfg.Logger.ErrorLog.Fatal(err)
}
