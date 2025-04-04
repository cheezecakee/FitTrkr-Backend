package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	FirstName    string
	LastName     string
	PasswordHash string
	Email        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsPremium    bool
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsPremium bool      `json:"is_premium"`
}
