package usecase

import (
	"context"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// CreateUser handles user creation from Kafka events
func (uc *UserUsecase) CreateUser(ctx context.Context, event events.UserCreatedEvent) error {
	user := &domain.User{
		ID:          event.ID,
		Name:        event.Name,
		Email:       event.Email,
		PhoneNumber: event.PhoneNumber,
		Role:        event.Role,
		CreatedAt:   event.CreatedAt,
	}

	return uc.userRepo.CreateUser(ctx, user)
}

// UpdateUser handles user updates from Kafka events
func (uc *UserUsecase) UpdateUser(ctx context.Context, event events.UserUpdatedEvent) error {
	user := &domain.User{
		ID:          event.ID,
		Name:        event.Name,
		Email:       event.Email,
		PhoneNumber: event.PhoneNumber,
		Role:        event.Role,
		UpdatedAt:   event.UpdatedAt,
	}

	return uc.userRepo.UpdateUser(ctx, user)
}
