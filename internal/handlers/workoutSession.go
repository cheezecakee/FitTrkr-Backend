package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/cheezecakee/FitLogr/internal/services"
)

// GET
// Directs the user either to sessions/workouts or sessions/start/workouts/{id}
func (cfg *Config) GetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get userID from context
	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	fmt.Println("userID: ", userID)

	// Get the current workout session from Redis
	session, err := cfg.DB.GetWorkoutSession(ctx, userID)
	if err != nil {
		cfg.Logger.InfoLog.Println("Workout session already in progress!")
		cfg.Helper.ServerError(w, err)
		// Ensure proper redirect with the correct status code
		json.NewEncoder(w).Encode(session.ID)
		return
	}

	// Redirect to the active session
	json.NewEncoder(w).Encode(map[string]any{
		"message": "No session found! Start a new session!",
	})
}

// GET
// If no session availble lists availble workouts to start a session
func (cfg *Config) ListSessionWorkouts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	sessionData, err := cfg.DB.GetWorkoutsByID(ctx, userID)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sessionData); err != nil {
		cfg.Helper.ServerError(w, err)
	}
}

// POST
// Double tap to start a session
func (cfg *Config) StartSessionWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workoutIDStr := chi.URLParam(r, "id")
	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	// Extract user ID from context
	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	// Get exercise IDs from DB
	totalExercises, err := cfg.DB.GetTotalExercises(ctx, int32(workoutID))
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	params := services.CreateWorkoutSessionParams{
		UserID:         userID,
		WorkoutID:      int32(workoutID),
		TotalExercises: int32(totalExercises),
	}
	session, err := cfg.DB.CreateWorkoutSession(ctx, params)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Workout session created!")

	w.Header().Set("Content-Type", "application/json")

	// return the new session id
	json.NewEncoder(w).Encode(map[string]string{
		"message":             "Workout session started",
		"exercise_session_id": fmt.Sprintf("%d", session),
	})
}

// GET
// Lists exercises to start an exercise session from the current workout session
func (cfg *Config) ListSessionExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionIDStr := chi.URLParam(r, "sessionID")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	// Get exercise list from current workout in session
	exercises, err := cfg.DB.GetWorkoutExercises(ctx, int32(sessionID))
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
// Double tap starts an exercise session
func (cfg *Config) StartSessionExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exerciseIDStr := chi.URLParam(r, "id")
	exerciseID, err := strconv.Atoi(exerciseIDStr)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	sessionIDStr := chi.URLParam(r, "sessionID")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	// Get total sets from DB
	exercise, err := cfg.DB.GetWorkoutExerciseByID(ctx, services.GetWorkoutExerciseByIDParams{
		ID:   int32(exerciseID),
		ID_2: int32(sessionID),
	})
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	params := services.CreateExerciseSessionParams{
		WorkoutSessionID:  int32(sessionID),
		WorkoutExerciseID: int32(exerciseID),
		TotalSets:         exercise.Sets,
	}
	// Start exercise session in Redis
	session, err := cfg.DB.CreateExerciseSession(ctx, params)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	cfg.Logger.InfoLog.Println("Exercise session started!")
	json.NewEncoder(w).Encode(map[string]string{
		"message":             "Exercise session started",
		"exercise_session_id": fmt.Sprintf("%d", session),
	})
}

func (cfg *Config) LogSet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract user ID from context
	userID, ok := r.Context().Value(cfg.ContextKey["userID"]).(uuid.UUID)
	if !ok {
		cfg.Helper.InfoLog.Println("From StartSessionExercise")
		cfg.Helper.ClientError(w, http.StatusUnauthorized)
		return
	}

	// Extract session IDs from path
	workoutSessionIDStr := chi.URLParam(r, "sessionID")
	workoutSessionID, err := strconv.Atoi(workoutSessionIDStr)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	exerciSessionIDStr := chi.URLParam(r, "exerciseSessionID")
	exerciseSessionID, err := strconv.Atoi(exerciSessionIDStr)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	data := json.NewDecoder(r.Body)
    if err := data.Decode(v any); err != nil {
        cfg.Helper.ServerError(w, err)
    }
}

// Single tap workout to view exercises avialable inside
func (cfg *Config) ViewSessionWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workoutIDStr := chi.URLParam(r, "id")
	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	// Fetch exercises for this workout from DB
	exercises, err := cfg.DB.GetWorkoutExercises(ctx, int32(workoutID))
	if err != nil {
		cfg.Helper.ServerError(w, err)
		return
	}

	// Send JSON response with exercises
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(exercises); err != nil {
		cfg.Helper.ServerError(w, err)
	}
}

// Single tap exercise to view exercise details
func (cfg *Config) ViewSessionExerciseDetails(w http.ResponseWriter, r *http.Request) {
}
