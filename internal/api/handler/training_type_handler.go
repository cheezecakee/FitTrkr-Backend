package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cheezecakee/fitrkr/internal/db/exercise"
)

type TrainingTypeHandler struct {
	service exercise.TrainingTypeService
}

func NewTrainingTypeHandler(service exercise.TrainingTypeService) *TrainingTypeHandler {
	return &TrainingTypeHandler{service: service}
}

// List lists all exercise types
// @Summary List all exercise types
// @Tags exercise-types
// @Security BearerAuth
// @Produce json
// @Success 200 {array} exercise.TrainingType
// @Router /api/v1/admin/exercise-types [get]
func (h *TrainingTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 1000 // or any reasonable default
	types, err := h.service.List(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, "Failed to list exercise types", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"data":  types,
		"error": nil,
	})
}
