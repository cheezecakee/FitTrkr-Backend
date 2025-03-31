package config

import (
	"log"
	"os"
)

type Config struct {
	DBConnString string
	Port         string
}

func LoadConfig() Config {
	dbConn := os.Getenv("DB_CONN_STRING")
	if dbConn == "" {
		log.Fatal("DB_CONN_STRING environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		DBConnString: dbConn,
		Port:         port,
	}
}
