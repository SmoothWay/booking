package handler

import (
	"net/http"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/validator"
	"github.com/SmoothWay/booking/internal/util"
)

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {

	var req PageRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	bookings, err := h.bookingUsecase.GetBookings(r.Context(), req.Page, req.PageSize)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, GetBookingsResponse{
		Pageable: Pageable{
			Total:    len(bookings),
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Bookings: bookingsToResponse(bookings),
	})
}

func bookingsToResponse(bookings []domain.Booking) []Booking {
	response := make([]Booking, len(bookings))
	for i, booking := range bookings {
		response[i] = Booking{
			ID:        booking.ID,
			UserID:    booking.UserID,
			UnitID:    booking.UnitID,
			StartDate: booking.StartDate,
			EndDate:   booking.EndDate,
		}
	}
	return response
}
