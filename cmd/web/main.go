package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/cheezecakee/FitLogr/internal/database"
)

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
	dbQueries := database.New(db)

	// Initialize jwt secret
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if dbURL == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	apiCfg := &ApiConfig{
		DB:        dbQueries,
		JWTSecret: jwtSecret,
	}

	// Setup routes
	router := apiCfg.SetupRoutes(dbQueries) // Pass dbQueries to your router

	// Start server
	port := ":8080"

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Println("Server is running on http://localhost" + port)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	log.Fatal(err)
}
