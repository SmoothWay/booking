package handler

import (
	"log"
	"net/http"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/validator"
	"github.com/SmoothWay/booking/internal/util"
)

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {

	var req PageRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 10
	}

	if err := validator.Validate(req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	bookings, err := h.bookingUsecase.GetBookings(r.Context(), req.Page, req.PageSize)
	if err != nil {
		log.Println("getbookings error:", err)
		util.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrInternalServerError,
		})
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

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := validator.Validate(req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	err := h.bookingUsecase.CreateBooking(r.Context(), events.BookingCreatedEvent{
		UserID:    req.UserID,
		UnitID:    req.UnitID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})

	if err != nil {
		log.Println("createbooking error:", err)
		util.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrInternalServerError,
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, CreateBookingResponse{
		ID: req.UserID,
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
