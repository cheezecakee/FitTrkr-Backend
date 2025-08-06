package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr/internal/db/exercise"
)

// ExerciseCategoryHandler handles HTTP requests for exercise category-related operations
type ExerciseCategoryHandler struct {
	service exercise.CategoryService
}

// NewExerciseCategoryHandler creates a new ExerciseCategoryHandler
func NewExerciseCategoryHandler(service exercise.CategoryService) *ExerciseCategoryHandler {
	return &ExerciseCategoryHandler{service: service}
}

// List lists all exercise categories
// @Summary List all exercise categories
// @Tags exercise-categories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} exercise.Category
// @Router /api/v1/admin/exercise-categories [get]
func (h *ExerciseCategoryHandler) List(w http.ResponseWriter, r *http.Request) {
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
	categories, err := h.service.List(ctx, offset, limit)
	if err != nil {
		if err.Error() == "limit must be greater than 0" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to list exercise categories", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  categories,
		"error": nil,
	})
}

// GetByID gets an exercise category by ID
// @Summary Get exercise category by ID
// @Tags exercise-categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "Exercise Category ID"
// @Success 200 {object} exercise.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/admin/exercise-categories/{id} [get]
func (h *ExerciseCategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exercise category ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	category, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "valid exercise category ID is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err.Error() == "exercise category not found" {
			http.Error(w, "Exercise category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get exercise category", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  category,
		"error": nil,
	})
}
