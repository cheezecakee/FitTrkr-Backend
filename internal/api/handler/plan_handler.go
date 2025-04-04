package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/cheezecakee/fitrkr/internal/models"
	"github.com/cheezecakee/fitrkr/internal/service"
)

type PlanHandler struct {
	svc service.PlanService
}

func NewPlanHandler(svc service.PlanService) *PlanHandler {
	return &PlanHandler{svc: svc}
}

func (h *PlanHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.svc.List(r.Context(), 0, 10)
	if err != nil {
		ServerError(w, err)
		return
	}
	Response(w, http.StatusOK, plans)
}

func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	var plan models.Plan
	if err := json.NewDecoder(r.Body); err != nil {
		ClientError(w, http.StatusBadRequest)
		return
	}
	created, err := h.svc.Create(r.Context(), &plan)
	if err != nil {
		if err == service.ErrInvalidPlan {
			ClientError(w, http.StatusBadRequest)
			return
		} else {
			ServerError(w, err)
		}
		return
	}
	Response(w, http.StatusCreated, created)
}

func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "planID")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ClientError(w, http.StatusBadRequest)
		return
	}

	plan, err := h.svc.GetByID(r.Context(), uint(id))
	if err != nil {
		if err == service.ErrPlanNotFound {
			ClientError(w, http.StatusBadRequest)
			return
		} else {
			ServerError(w, err)
		}
		return
	}
	Response(w, http.StatusOK, plan)
}

func (h *PlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
}

func (h *PlanHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {}
