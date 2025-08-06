package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr/internal/db/exercise"
)

// EquipmentHandler handles HTTP requests for equipment-related operations
type EquipmentHandler struct {
	service exercise.EquipmentService
}

// NewEquipmentHandler creates a new EquipmentHandler
func NewEquipmentHandler(service exercise.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{service: service}
}

// List lists all equipment
// @Summary List all equipment
// @Tags equipment
// @Security BearerAuth
// @Produce json
// @Success 200 {array} exercise.Equipment
// @Router /api/v1/admin/equipment [get]
func (h *EquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
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
	equipment, err := h.service.List(ctx, offset, limit)
	if err != nil {
		if err.Error() == "limit must be greater than 0" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to list equipment", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  equipment,
		"error": nil,
	})
}

// GetByID gets equipment by ID
// @Summary Get equipment by ID
// @Tags equipment
// @Security BearerAuth
// @Produce json
// @Param id path int true "Equipment ID"
// @Success 200 {object} exercise.Equipment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/admin/equipment/{id} [get]
func (h *EquipmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	equipment, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "valid equipment ID is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err.Error() == "equipment not found" {
			http.Error(w, "Equipment not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get equipment", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":  equipment,
		"error": nil,
	})
} 