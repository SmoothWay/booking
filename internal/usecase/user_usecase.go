package usecase

import (
	"context"
	"encoding/json"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/broker/kafka"
	"github.com/google/uuid"
)

type UserUsecase struct {
	kafkaProducer *kafka.Producer
	userRepo      domain.UserRepository
}

func NewUserUsecase(kafkaProducer *kafka.Producer, userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		kafkaProducer: kafkaProducer,
		userRepo:      userRepo,
	}
}

// CreateUser handles user creation from Kafka events
func (uc *UserUsecase) CreateUser(ctx context.Context, event events.UserCreatedEvent) error {
	user := &domain.User{
		ID:          uuid.New().String(),
		Name:        event.Name,
		Email:       event.Email,
		PhoneNumber: event.PhoneNumber,
		Role:        event.Role,
		Password:    event.Password,
		CreatedAt:   event.CreatedAt,
	}
	err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	json, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = uc.kafkaProducer.SendMessage(ctx, "user_created", json)
	if err != nil {
		return err
	}

	return nil
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

func (uc *UserUsecase) GetUsers(ctx context.Context, page int, pageSize int) ([]domain.User, error) {
	users, err := uc.userRepo.GetUsers(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var usersDomain []domain.User
	for _, user := range users {
		usersDomain = append(usersDomain, *user)
	}
	return usersDomain, nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}
