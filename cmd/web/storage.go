package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	dataFile = "./internal/data.json"
	mu       sync.Mutex
)

// Data structure for JSON storage
type Data struct {
	Users     []User      `json:"users"`
	Workouts  []Workout   `json:"workouts"`
	Exercises []Exercises `json:"exercises"`
}

type User struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Age          int32     `json:"age"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Workout struct {
	ID              string            `json:"id"`
	UserID          string            `json:"user_id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	WorkoutExercise []WorkoutExercise `json:"workout_exercise"`
	CreatedAt       time.Time         `json:"createdat"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type WorkoutExercise struct {
	ID         string    `json:"id"`
	WorkoutID  string    `json:"workout_id"`
	ExerciseID string    `json:"exercise_id"`
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

type Exercises struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkoutSession struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	WorkoutID  string     `json:"workout_id"`
	IsFinished bool       `json:"is_finished"`
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type WorkoutExerciseSession struct {
	ID                string    `json:"id"`
	WorkoutSessionID  string    `json:"workout_session_id"`
	WorkoutExerciseID string    `json:"workout_exercise_id"`
	IsFinished        bool      `json:"is_finished"`
	Skipped           bool      `json:"skipped"`
	StartedAt         time.Time `json:"started_at"`
	EndedAt           time.Time `json:"ended_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type WorkoutExerciseLogs struct {
	ID                string    `json:"id"`
	WorkoutExerciseID string    `json:"workout_exercise_id"`
	SetNumber         int       `json:"set_number"`
	Reps              int       `json:"reps"`
	Weight            float64   `json:"weight"`
	Interval          int64     `json:"interval"`
	Notes             string    `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
}

type SessionLogs struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"user_id"`
	WorkoutSessionID      string    `json:"workout_session_id"`
	WorkoutExerciseLogsID *string   `json:"workout_exercise_logs_id"`
	LogType               string    `json:"log_type"`
	LogPriority           string    `json:"log_priority"`
	LogMessage            string    `json:"log_message"`
	CreatedAt             time.Time `json:"created_at"`
}

// Load data from JSON file
func LoadData(dataFile string) (*Data, error) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile(dataFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return &Data{}, err
	}

	var data Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return &Data{}, err
	}

	return &data, nil
}

// Save data back to JSON file
func SaveData(data *Data) error {
	mu.Lock()
	defer mu.Unlock()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	err = os.WriteFile(dataFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return err
	}

	return nil
}
