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

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {

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

	bookings, err := h.bookingUsecase.GetBookings(r.Context(), page, pageSize)
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
			Page:     page,
			PageSize: pageSize,
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
