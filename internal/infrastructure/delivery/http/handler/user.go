package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/validator"
	"github.com/SmoothWay/booking/internal/util"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

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

	users, err := h.userUsecase.GetUsers(r.Context(), page, pageSize)
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
			Page:     page,
			PageSize: pageSize,
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
		if err == domain.ErrUserAlreadyExists {
			util.WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "User already exists",
			})
			return
		}
		util.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: ErrInternalServerError,
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, nil)
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
