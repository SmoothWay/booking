package handler

import (
	"net/http"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/validator"
	"github.com/SmoothWay/booking/internal/util"
)

func (h *Handler) GetUnits(w http.ResponseWriter, r *http.Request) {
	var req PageRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	units, err := h.unitUsecase.GetUnits(r.Context(), req.Page, req.PageSize)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, GetUnitsResponse{
		Pageable: Pageable{
			Total:    len(units),
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Units: unitsToResponse(units),
	})
}

func (h *Handler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	var req CreateUnitRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.unitUsecase.CreateUnit(r.Context(), events.UnitCreatedEvent{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
	})
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, nil)
}

func unitsToResponse(units []domain.Unit) []Unit {
	response := make([]Unit, len(units))
	for i, unit := range units {
		response[i] = Unit{
			ID:   unit.ID,
			Name: unit.Name,
		}
	}
	return response
}
