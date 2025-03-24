package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID        int64
	UserID    uuid.UUID
	PlanID    uint
	Type      string
	Priority  string
	Message   string
	Pr        bool
	Context   string // JSON string
	Timestamp time.Time
}
