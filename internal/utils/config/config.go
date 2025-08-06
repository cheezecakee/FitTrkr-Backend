// Package config provides configuration utilities for FitTrkr.
package config

import (
	"log"
	"os"
	"time"

	"github.com/cheezecakee/fitrkr/internal/utils/auth"
)

type Config struct {
	DBConnString string
	Port         string
	JWTManager   auth.JWT
}

func LoadConfig() Config {
	dbConn := os.Getenv("DB_CONN_STRING")
	if dbConn == "" {
		log.Fatal("DB_CONN_STRING environment variable is required")
	}
	log.Println("dbConn: ", dbConn)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if string(jwtSecret) == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	// 15 day expiration
	jwtManager := auth.NewJWTManager(jwtSecret, time.Hour*360)

	return Config{
		DBConnString: dbConn,
		Port:         port,
		JWTManager:   jwtManager,
	}
}
