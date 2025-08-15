package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetUsers(ctx context.Context, page int, pageSize int) ([]*domain.User, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) SendMessage(ctx context.Context, topic string, message []byte) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func TestUserUsecase_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockKafkaProducer := new(MockKafkaProducer)
	usecase := NewUserUsecase(mockKafkaProducer, mockRepo)

	event := events.UserCreatedEvent{
		ID:          "user-123",
		Name:        "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "1234567890",
		Role:        "user",
		Password:    "hashedpassword",
		CreatedAt:   time.Now(),
	}
	mockKafkaProducer.On("SendMessage", mock.Anything, "user.created", mock.Anything).Return(nil)
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == event.ID && user.Email == event.Email
	})).Return(nil)

	err := usecase.CreateUser(context.Background(), event)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_CreateUser_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockKafkaProducer := new(MockKafkaProducer)
	usecase := NewUserUsecase(mockKafkaProducer, mockRepo)

	event := events.UserCreatedEvent{
		ID:          "user-123",
		Name:        "John Doe",
		Email:       "existing@example.com",
		PhoneNumber: "1234567890",
		Role:        "user",
		Password:    "hashedpassword",
		CreatedAt:   time.Now(),
	}

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(domain.ErrUserAlreadyExists)

	err := usecase.CreateUser(context.Background(), event)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_GetUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockKafkaProducer := new(MockKafkaProducer)
	usecase := NewUserUsecase(mockKafkaProducer, mockRepo)

	expectedUsers := []*domain.User{
		{
			ID:          "user-1",
			Name:        "User 1",
			Email:       "user1@example.com",
			PhoneNumber: "1234567890",
			Role:        "user",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "user-2",
			Name:        "User 2",
			Email:       "user2@example.com",
			PhoneNumber: "0987654321",
			Role:        "admin",
			CreatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetUsers", mock.Anything, 1, 10).Return(expectedUsers, nil)

	result, err := usecase.GetUsers(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedUsers[0].ID, result[0].ID)
	assert.Equal(t, expectedUsers[1].ID, result[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_UpdateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockKafkaProducer := new(MockKafkaProducer)
	usecase := NewUserUsecase(mockKafkaProducer, mockRepo)

	event := events.UserUpdatedEvent{
		ID:          "user-123",
		Name:        "Updated Name",
		Email:       "updated@example.com",
		PhoneNumber: "9876543210",
		Role:        "admin",
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == event.ID && user.Name == event.Name
	})).Return(nil)

	err := usecase.UpdateUser(context.Background(), event)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_UpdateUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockKafkaProducer := new(MockKafkaProducer)
	usecase := NewUserUsecase(mockKafkaProducer, mockRepo)

	event := events.UserUpdatedEvent{
		ID:          "non-existent",
		Name:        "Updated Name",
		Email:       "updated@example.com",
		PhoneNumber: "9876543210",
		Role:        "admin",
		UpdatedAt:   time.Now(),
	}

	expectedError := errors.New("user not found")
	mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(expectedError)

	err := usecase.UpdateUser(context.Background(), event)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// Test table-driven tests for different user roles
func TestUserUsecase_UserRoleValidation(t *testing.T) {
	tests := []struct {
		name      string
		role      string
		isValid   bool
		setupMock func(*MockUserRepository)
	}{
		{
			name:    "Valid User Role",
			role:    "user",
			isValid: true,
			setupMock: func(mockRepo *MockUserRepository) {
				user := &domain.User{ID: "1", Role: "user"}
				mockRepo.On("GetUserByID", mock.Anything, "1").Return(user, nil)
			},
		},
		{
			name:    "Valid Admin Role",
			role:    "admin",
			isValid: true,
			setupMock: func(mockRepo *MockUserRepository) {
				user := &domain.User{ID: "1", Role: "admin"}
				mockRepo.On("GetUserByID", mock.Anything, "1").Return(user, nil)
			},
		},
		{
			name:    "Valid Guest Role",
			role:    "guest",
			isValid: true,
			setupMock: func(mockRepo *MockUserRepository) {
				user := &domain.User{ID: "1", Role: "guest"}
				mockRepo.On("GetUserByID", mock.Anything, "1").Return(user, nil)
			},
		},
		{
			name:    "Invalid Role",
			role:    "invalid",
			isValid: false,
			setupMock: func(mockRepo *MockUserRepository) {
				user := &domain.User{ID: "1", Role: "invalid"}
				mockRepo.On("GetUserByID", mock.Anything, "1").Return(user, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockKafkaProducer := new(MockKafkaProducer)
			usecase := NewUserUsecase(mockKafkaProducer, mockRepo)

			tt.setupMock(mockRepo)

			user, err := usecase.GetUserByID(context.Background(), "1")

			assert.NoError(t, err)
			assert.Equal(t, tt.role, user.Role)

			// Test role-specific methods
			switch tt.role {
			case "admin":
				assert.True(t, user.IsAdmin())
				assert.False(t, user.IsUser())
				assert.False(t, user.IsGuest())
			case "user":
				assert.False(t, user.IsAdmin())
				assert.True(t, user.IsUser())
				assert.False(t, user.IsGuest())
			case "guest":
				assert.False(t, user.IsAdmin())
				assert.False(t, user.IsUser())
				assert.True(t, user.IsGuest())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Test context cancellation
func TestUserUsecase_ContextCancellation(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockKafkaProducer := new(MockKafkaProducer)
	usecase := NewUserUsecase(mockKafkaProducer, mockRepo)

	mockRepo.On("GetUsers", mock.Anything, 1, 10).Return(nil, context.Canceled)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := usecase.GetUsers(ctx, 1, 10)

	assert.Error(t, err)
}

// Test concurrent access
func TestUserUsecase_ConcurrentAccess(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockKafkaProducer := new(MockKafkaProducer)
	usecase := NewUserUsecase(mockKafkaProducer, mockRepo)

	user := &domain.User{ID: "1", Name: "Test User"}
	mockRepo.On("GetUserByID", mock.Anything, "1").Return(user, nil).Times(10)

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := usecase.GetUserByID(context.Background(), "1")
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	mockRepo.AssertExpectations(t)
}
