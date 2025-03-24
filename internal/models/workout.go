package models

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	ID              int               `json:"id"`
	UserID          uuid.UUID         `json:"user_id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	WorkoutExercise []WorkoutExercise `json:"workout_exercise"`
	CreatedAt       time.Time         `json:"createdat"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
