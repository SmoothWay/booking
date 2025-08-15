package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/validator"
	"github.com/SmoothWay/booking/internal/util"
)

func (h *Handler) GetUnits(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	if err := validator.Validate(PageRequest{
		Page:     page,
		PageSize: pageSize,
	}); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	units, err := h.unitUsecase.GetUnits(r.Context(), page, pageSize)
	if err != nil {
		log.Println("getunits error:", err)
		util.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrInternalServerError,
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, GetUnitsResponse{
		Pageable: Pageable{
			Total:    len(units),
			Page:     page,
			PageSize: pageSize,
		},
		Units: unitsToResponse(units),
	})
}

func (h *Handler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	var req CreateUnitRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	err := h.unitUsecase.CreateUnit(r.Context(), events.UnitCreatedEvent{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
	})
	if err != nil {
		log.Println("createunit error:", err)
		util.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrInternalServerError,
		})
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
