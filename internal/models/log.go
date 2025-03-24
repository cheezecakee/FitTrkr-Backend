package models

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type SessionLog struct {
	ExerciseLogs []ExerciseLog `json:"exercise_logs"`
	WorkoutLog   WorkoutLog    `json:"session_log"`
}

type WorkoutLog struct {
	ID             big.Int     `json:"id"`
	UserID         uuid.UUID   `json:"user_id"`
	WorkoutId      int         `json:"workout_id"`
	ExerciseIDs    []int       `json:"exercise_ids"`
	TotalExercises int         `json:"total_exercises"`
	Completed      bool        `json:"complete"`
	Log            string      `json:"log"`
	LogType        string      `json:"log_type"`
	LogPriority    string      `json:"log_priority"`
	LogMessage     string      `json:"log_message"`
	PRs            map[int]any `json:"prs"` // Map of exercise ID to weight or duration(cardio)
	PR             bool        `json:"pr"`
	Notes          string      `json:"notes"`
	CreatedAt      time.Time   `json:"created_at"`
	Timestamp      int64       `json:"timestamp"` // Unix timestamp
	Duration       int         `json:"duration"`
}

type ExerciseLog struct {
	ExerciseID   int32   `json:"exercise_ids"`
	ExerciseName string  `json:"exercise_name"`
	Reps         int     `json:"reps"`
	Sets         int     `json:"sets"`
	Weight       float64 `json:"weight"`
	Interval     int     `json:"interval"` // time in seconds
	Rest         int     `json:"rest"`     // time in seconds
	Duration     int     `json:"duration"`
	Completed    bool    `json:"complete"`
	Notes        string  `json:"notes"`
	Timestamp    int64   `json:"timestamp"` // Unix timestamp
}
