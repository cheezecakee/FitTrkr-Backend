package handler

import (
	"net/http"

	"github.com/cheezecakee/fitrkr/internal/service"
)

type PlanExHandler struct {
	svc service.PlanExService
}

func NewPlanExHandler(svc service.PlanExService) *PlanExHandler {
	return &PlanExHandler{svc: svc}
}

func (h *PlanExHandler) GetPlanExercises(w http.ResponseWriter, r *http.Request) {}

func (h *PlanExHandler) CreatePlanExercise(w http.ResponseWriter, r *http.Request) {}

func (h *PlanExHandler) UpdatePlanExercise(w http.ResponseWriter, r *http.Request) {}

func (h *PlanExHandler) DeletePlanExercise(w http.ResponseWriter, r *http.Request) {}
