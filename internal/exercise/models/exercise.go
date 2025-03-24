package models

import (
	"time"
)

type Exercise struct {
	ID          uint
	Name        string
	Description string
	Category    string // e.g., "Chest", "Legs", etc.
	Equipment   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
