package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/cheezecakee/FitLogr/internal/models"
	"github.com/cheezecakee/FitLogr/internal/services"
)

// GET
func (cfg *Config) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	workouts, err := cfg.DB.GetWorkoutsByID(ctx, userID)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workouts); err != nil {
		cfg.Helper.ServerError(w, err)
	}
}

// GET
func (cfg *Config) GetWorkoutsExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workoutIDStr := r.PathValue("id")

	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	exercises, err := cfg.DB.GetWorkoutExercises(ctx, int32(workoutID))
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(exercises); err != nil {
		cfg.Helper.ServerError(w, err)
	}
}

// POST
func (cfg *Config) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	var req models.Workout
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	workout, err := cfg.DB.CreateWorkout(ctx, services.CreateWorkoutParams{
		UserID:      userID,
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	})
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Workout created successfully!")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// PUT
func (cfg *Config) EditWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workoutIDStr := r.PathValue("id")

	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	updatedWorkout, err := cfg.DB.EditWorkout(ctx, services.EditWorkoutParams{
		ID:          int32(workoutID),
		UserID:      userID,
		Name:        cfg.Helper.StringToNullString(req.Name),
		Description: cfg.Helper.StringToNullString(req.Description),
	})
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Workout updated successfully!")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedWorkout)
}

// DELETE
func (cfg *Config) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workoutIDStr := r.PathValue("id")

	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	if err := cfg.DB.DeleteWorkout(ctx, int32(workoutID)); err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Workout deleted successfully!")
	w.WriteHeader(http.StatusNoContent)
}

// POST
func (cfg *Config) CreateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.WorkoutExercise
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	workoutIDStr := r.PathValue("id")

	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	totalRestSeconds := (req.RestMin * 60) + req.RestSec
	log.Println("workoutID: ", req.WorkoutID)

	workoutExercise, err := cfg.DB.CreateWorkoutExercise(ctx, services.CreateWorkoutExerciseParams{
		WorkoutID:  int32(workoutID),
		ExerciseID: int32(req.ExerciseID),
		Sets:       req.Sets,
		RepsMin:    req.RepsMin,
		RepsMax:    req.RepsMax,
		Weight:     req.Weight,
		Interval:   req.Interval,
		Rest:       totalRestSeconds,
	})
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Exercise added to workout successfully!")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workoutExercise)
}

// PUT
func (cfg *Config) EditWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exerciseIDStr := r.PathValue("id")
	exerciseID, err := strconv.Atoi(exerciseIDStr)
	if err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	var req struct {
		Sets     int32   `json:"sets"`
		RepsMin  int32   `json:"reps_min"`
		RepsMax  int32   `json:"reps_max"`
		Weight   float64 `json:"weight"`
		Interval int32   `json:"interval"`
		Rest     int32   `json:"rest"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	updatedExercise, err := cfg.DB.EditWorkoutExercise(ctx, services.EditWorkoutExerciseParams{
		ID:       int32(exerciseID),
		Sets:     sql.NullInt32{Int32: req.Sets, Valid: true},
		RepsMin:  sql.NullInt32{Int32: req.RepsMin, Valid: true},
		RepsMax:  sql.NullInt32{Int32: req.RepsMax, Valid: true},
		Weight:   sql.NullFloat64{Float64: req.Weight, Valid: true},
		Interval: sql.NullInt32{Int32: req.Interval, Valid: true},
		Rest:     sql.NullInt32{Int32: req.Rest, Valid: true},
	})
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Workout exercise updated successfully!")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExercise)
}

// DELETE
func (cfg *Config) DeleteWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exerciseIDStr := r.PathValue("id")
	exerciseID, err := strconv.Atoi(exerciseIDStr)
	if err != nil {
		cfg.Helper.ClientError(w, http.StatusBadRequest)
		return
	}

	if err := cfg.DB.DeleteWorkoutExercise(ctx, int32(exerciseID)); err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Workout Exercise deleted successfully!")
	w.WriteHeader(http.StatusNoContent)
}
