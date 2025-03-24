package models

import (
	"time"
)

type WorkoutExercise struct {
	ID         int32     `json:"id"`
	WorkoutID  int32     `json:"workout_id"`
	ExerciseID int32     `json:"exercise_id"`
	RepsMin    int32     `json:"reps_min"`
	RepsMax    int32     `json:"reps_max"`
	Weight     float64   `json:"weight"`
	Sets       int32     `json:"sets"`
	Interval   int32     `json:"interval"`
	RestMin    int32     `json:"rest_min"`
	RestSec    int32     `json:"rest_sec"`
	CreatedAt  time.Time `json:"createdat"`
	UpdatedAt  time.Time `json:"updated_at"`
}
