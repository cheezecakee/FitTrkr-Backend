package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PlanID    uint
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Exercises []SessionEx
}

type SessionEx struct {
	ID         uuid.UUID
	SessionID  uuid.UUID
	ExerciseID uint
	Order      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Sets       []ExSet
	Exercise   Exercise
}

type ExSet struct {
	ID                uint
	SessionExerciseID uuid.UUID
	Number            int
	Reps              int
	Weight            float64
	Duration          int
	Distance          float64 // For cardio
	Notes             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
