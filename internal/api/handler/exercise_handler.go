package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr/internal/db/exercise"
)

// ExerciseHandler handles HTTP requests for exercise-related operations
type ExerciseHandler struct {
	service exercise.ExerciseService
}

// NewExerciseHandler creates a new ExerciseHandler
func NewExerciseHandler(service exercise.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{service: service}
}

// Create creates a new exercise
// @Summary Create a new exercise
// @Tags exercises
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body exercise.CreateExerciseRequest true "Exercise payload"
// @Success 201 {object} exercise.Exercise
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/exercises [post]
func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req exercise.CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	created, err := h.service.Create(ctx, &req)
	if err != nil {
		switch err.Error() {
		case "exercise already exists", "exercise name is required", "exercise name must not exceed 100 characters",
			"exercise description is required", "valid category ID is required",
			"valid equipment ID is required", "exercise type must be 'strength' or 'cardio'":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Failed to create exercise", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  created,
		"error": nil,
	})
}

// GetByID gets an exercise by ID
// @Summary Get exercise by ID
// @Tags exercises
// @Security BearerAuth
// @Produce json
// @Param id path int true "Exercise ID"
// @Success 200 {object} exercise.Exercise
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/admin/exercises/{id} [get]
func (h *ExerciseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ex, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "valid exercise ID is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err.Error() == "exercise not found" {
			http.Error(w, "Exercise not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get exercise", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  ex,
		"error": nil,
	})
}

// List lists all exercises
// @Summary List all exercises
// @Tags exercises
// @Security BearerAuth
// @Produce json
// @Success 200 {array} exercise.Exercise
// @Router /api/v1/admin/exercises [get]
func (h *ExerciseHandler) List(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	ctx := r.Context()
	exercises, err := h.service.List(ctx, offset, limit)
	if err != nil {
		if err.Error() == "limit must be greater than 0" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to list exercises", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  exercises,
		"error": nil,
	})
}

// Update updates an exercise
// @Summary Update an exercise
// @Tags exercises
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Exercise ID"
// @Param request body exercise.UpdateExerciseRequest true "Exercise payload"
// @Success 200 {object} exercise.UpdateExerciseRequest
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/admin/exercises/{id} [put]
func (h *ExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	var updateReq exercise.UpdateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	updated, err := h.service.Update(ctx, &updateReq, id)
	if err != nil {
		switch err.Error() {
		case "valid exercise ID is required", "exercise name is required",
			"exercise name must not exceed 100 characters", "exercise description is required",
			"valid category ID is required", "valid equipment ID is required",
			"exercise type must be 'strength' or 'cardio'":
			http.Error(w, err.Error(), http.StatusBadRequest)
		case "exercise not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  updated,
		"error": nil,
	})
}

// Delete deletes an exercise
// @Summary Delete an exercise
// @Tags exercises
// @Security BearerAuth
// @Produce json
// @Param id path int true "Exercise ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/admin/exercises/{id} [delete]
func (h *ExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.Delete(ctx, id); err != nil {
		if err.Error() == "valid exercise ID is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err.Error() == "exercise not found" {
			http.Error(w, "Exercise not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete exercise", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get exercise with full details (including relationships)
// @Summary Get exercise with full details
// @Tags exercises
// @Security BearerAuth
// @Produce json
// @Param id path int true "Exercise ID"
// @Success 200 {object} exercise.Exercise
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/admin/exercises/{id}/details [get]
func (h *ExerciseHandler) GetWithDetails(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ex, err := h.service.GetExerciseWithDetails(ctx, id)
	if err != nil {
		if err.Error() == "valid exercise ID is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to get exercise details", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  ex,
		"error": nil,
	})
}

// Get exercises by muscle group
// @Summary Get exercises by muscle group
// @Tags exercises
// @Security BearerAuth
// @Produce json
// @Param muscle_group query string true "Muscle group name"
// @Success 200 {array} exercise.Exercise
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/exercises/by-muscle [get]
func (h *ExerciseHandler) GetByMuscleGroup(w http.ResponseWriter, r *http.Request) {
	muscleGroup := r.URL.Query().Get("muscle_group")
	if muscleGroup == "" {
		http.Error(w, "muscle_group parameter is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exercises, err := h.service.GetByMuscleGroupName(ctx, muscleGroup)
	if err != nil {
		http.Error(w, "Failed to get exercises by muscle group", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  exercises,
		"error": nil,
	})
}

// Get exercises by training type
// @Summary Get exercises by training type
// @Tags exercises
// @Security BearerAuth
// @Produce json
// @Param training_type query string true "Training type name"
// @Success 200 {array} exercise.Exercise
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/exercises/by-training-type [get]
func (h *ExerciseHandler) GetByTrainingType(w http.ResponseWriter, r *http.Request) {
	trainingType := r.URL.Query().Get("training_type")
	if trainingType == "" {
		http.Error(w, "training_type parameter is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exercises, err := h.service.GetByTrainingTypeName(ctx, trainingType)
	if err != nil {
		http.Error(w, "Failed to get exercises by training type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  exercises,
		"error": nil,
	})
}

// Search exercises (enhanced version that might include relationships in future)
// @Summary Search exercises
// @Tags exercises
// @Security BearerAuth
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} exercise.Exercise
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/exercises/search [get]
func (h *ExerciseHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "search query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exercises, err := h.service.Search(ctx, query)
	if err != nil {
		http.Error(w, "Failed to search exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  exercises,
		"error": nil,
	})
}
