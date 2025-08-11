package handler

import (
	"net/http"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/util"
)

type Handler struct {
	bookingUsecase domain.BookingUsecase
	userUsecase    domain.UserUsecase
	unitUsecase    domain.UnitUsecase
}

type PageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type Pageable struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func NewHandler(bookingUsecase domain.BookingUsecase, userUsecase domain.UserUsecase, unitUsecase domain.UnitUsecase) *Handler {
	return &Handler{bookingUsecase: bookingUsecase, userUsecase: userUsecase, unitUsecase: unitUsecase}
}

// TODO: implement liveness, readiness and health checks
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, "OK")
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, "OK")
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, "OK")
}
