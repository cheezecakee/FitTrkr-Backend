package models

import "time"

type Exercises struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updated_at"`
}
