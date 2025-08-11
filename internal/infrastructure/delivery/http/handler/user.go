package handler

import (
	"net/http"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http/validator"
	"github.com/SmoothWay/booking/internal/util"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req PageRequest
	if err := util.ReadJSON(r, &req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.userUsecase.GetUsers(r.Context(), req.Page, req.PageSize)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
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
