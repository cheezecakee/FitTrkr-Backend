package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cheezecakee/fitrkr/internal/db/exercise"
)

type MuscleGroupHandler struct {
	service exercise.MuscleGroupService
}

func NewMuscleGroupHandler(service exercise.MuscleGroupService) *MuscleGroupHandler {
	return &MuscleGroupHandler{service: service}
}

// List lists all muscle groups
// @Summary List all muscle groups
// @Tags muscle-groups
// @Security BearerAuth
// @Produce json
// @Success 200 {array} exercise.MuscleGroup
// @Router /api/v1/admin/muscle-groups [get]
func (h *MuscleGroupHandler) List(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 1000 // or any reasonable default
	groups, err := h.service.List(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, "Failed to list muscle groups", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"data":  groups,
		"error": nil,
	})
}
