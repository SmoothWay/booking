package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/validator"
	"github.com/SmoothWay/booking/internal/util"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.userUsecase.GetUsers(r.Context(), req.Page, req.PageSize)
	if err != nil {
		log.Println("getusers error:", err)
		util.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrInternalServerError,
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, GetUsersResponse{
		Pageable: Pageable{
			Total:    len(users),
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Users: usersToResponse(users),
	})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	err := h.userUsecase.CreateUser(r.Context(), events.UserCreatedEvent{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Role:        req.Role,
		Password:    req.Password,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		log.Println("createuser error:", err)
		util.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrInternalServerError,
		})
		return
	}
}

func usersToResponse(users []domain.User) []User {
	response := make([]User, len(users))
	for i, user := range users {
		response[i] = User{
			ID:   user.ID,
			Name: user.Name,
		}
	}
	return response
}
