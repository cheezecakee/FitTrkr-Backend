package repository

import (
	"database/sql"
	"log"
)

func NewDB(connString string) *sql.DB {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return db
}
