package models

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID          uint
	UserID      uuid.UUID
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Exercises   []PlanEx
}

type PlanEx struct {
	ID         uint
	PlanID     uint
	ExerciseID uint
	Name       string
	Order      int
	Sets       int
	Reps       [2]int
	Weight     float64
	Interval   int
	Rest       int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Exercise   Exercise
}
