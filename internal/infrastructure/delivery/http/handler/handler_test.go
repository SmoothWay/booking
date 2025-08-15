package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookingUsecase struct {
	mock.Mock
}

func (m *MockBookingUsecase) CreateBooking(ctx context.Context, event events.BookingCreatedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockBookingUsecase) GetBookings(ctx context.Context, page int, pageSize int) ([]domain.Booking, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Booking), args.Error(1)
}

func (m *MockBookingUsecase) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingUsecase) GetBookingsByUserID(ctx context.Context, userID string, page int, pageSize int) ([]domain.Booking, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Booking), args.Error(1)
}

func (m *MockBookingUsecase) CancelBooking(ctx context.Context, event events.BookingCancelledEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockBookingUsecase) UpdateBooking(ctx context.Context, event events.BookingUpdatedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) CreateUser(ctx context.Context, event events.UserCreatedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockUserUsecase) GetUsers(ctx context.Context, page int, pageSize int) ([]domain.User, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserUsecase) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1)
	}
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, event events.UserUpdatedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

type MockUnitUsecase struct {
	mock.Mock
}

func (m *MockUnitUsecase) CreateUnit(ctx context.Context, event events.UnitCreatedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockUnitUsecase) GetUnits(ctx context.Context, page int, pageSize int) ([]domain.Unit, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Unit), args.Error(1)
}

func (m *MockUnitUsecase) GetUnitByID(ctx context.Context, id string) (domain.Unit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.Unit{}, args.Error(1)
	}
	return args.Get(0).(domain.Unit), args.Error(1)
}

func TestHandler_CreateUser_Success(t *testing.T) {
	mockBookingUsecase := new(MockBookingUsecase)
	mockUserUsecase := new(MockUserUsecase)
	mockUnitUsecase := new(MockUnitUsecase)
	handler := NewHandler(mockBookingUsecase, mockUserUsecase, mockUnitUsecase)

	userData := map[string]any{
		"id":           "user-123",
		"name":         "John Doe",
		"email":        "john@example.com",
		"phone_number": "1234567890",
		"role":         "user",
		"password":     "hashedpassword",
		"created_at":   time.Now(),
	}

	requestBody, _ := json.Marshal(userData)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUserUsecase.On("CreateUser", mock.Anything, mock.MatchedBy(func(event events.UserCreatedEvent) bool {
		return event.Email == userData["email"] && event.Name == userData["name"]
	})).Return(nil)

	handler.CreateUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Empty(t, w.Body.String())

	mockUserUsecase.AssertExpectations(t)
}

func TestHandler_CreateUser_InvalidJSON(t *testing.T) {
	mockBookingUsecase := new(MockBookingUsecase)
	mockUserUsecase := new(MockUserUsecase)
	mockUnitUsecase := new(MockUnitUsecase)
	handler := NewHandler(mockBookingUsecase, mockUserUsecase, mockUnitUsecase)

	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "invalid")

	mockUserUsecase.AssertNotCalled(t, "CreateUser")
}

func TestHandler_CreateUser_DuplicateEmail(t *testing.T) {
	mockBookingUsecase := new(MockBookingUsecase)
	mockUserUsecase := new(MockUserUsecase)
	mockUnitUsecase := new(MockUnitUsecase)
	handler := NewHandler(mockBookingUsecase, mockUserUsecase, mockUnitUsecase)

	userData := map[string]any{
		"id":           "user-123",
		"name":         "John Doe",
		"email":        "existing@example.com",
		"phone_number": "1234567890",
		"role":         "user",
		"password":     "hashedpassword",
		"created_at":   time.Now(),
	}

	requestBody, _ := json.Marshal(userData)
	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUserUsecase.On("CreateUser", mock.Anything, mock.Anything).Return(domain.ErrUserAlreadyExists)

	handler.CreateUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "already exists")

	mockUserUsecase.AssertExpectations(t)
}

func TestHandler_GetUsers_Success(t *testing.T) {
	mockBookingUsecase := new(MockBookingUsecase)
	mockUserUsecase := new(MockUserUsecase)
	mockUnitUsecase := new(MockUnitUsecase)
	handler := NewHandler(mockBookingUsecase, mockUserUsecase, mockUnitUsecase)

	expectedUsers := []domain.User{
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

	req := httptest.NewRequest("GET", "/api/v1/user?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	mockUserUsecase.On("GetUsers", mock.Anything, 1, 10).Return(expectedUsers, nil)

	handler.GetUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["users"])

	mockUserUsecase.AssertExpectations(t)
}

func TestHandler_GetUsers_InvalidPagination(t *testing.T) {
	mockBookingUsecase := new(MockBookingUsecase)
	mockUserUsecase := new(MockUserUsecase)
	mockUnitUsecase := new(MockUnitUsecase)
	handler := NewHandler(mockBookingUsecase, mockUserUsecase, mockUnitUsecase)

	req := httptest.NewRequest("GET", "/api/v1/user?page=0&page_size=1001", nil)
	w := httptest.NewRecorder()

	mockUserUsecase.On("GetUsers", mock.Anything, 0, 1001).Return(nil, errors.New("invalid pagination"))

	handler.GetUsers(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// validator error message
	assert.Contains(t, response.Error, "cannot be")

	mockUserUsecase.AssertNotCalled(t, "GetUsers")
}

func TestHandler_GetUsers_DatabaseError(t *testing.T) {
	mockBookingUsecase := new(MockBookingUsecase)
	mockUserUsecase := new(MockUserUsecase)
	mockUnitUsecase := new(MockUnitUsecase)
	handler := NewHandler(mockBookingUsecase, mockUserUsecase, mockUnitUsecase)

	req := httptest.NewRequest("GET", "/api/v1/user?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	expectedError := errors.New("database connection failed")
	mockUserUsecase.On("GetUsers", mock.Anything, 1, 10).Return(nil, expectedError)

	handler.GetUsers(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInternalServerError, response.Error)

	mockUserUsecase.AssertExpectations(t)
}

func TestHandler_HTTPMethods(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		setupMocks     func(*MockBookingUsecase, *MockUserUsecase, *MockUnitUsecase)
	}{
		{
			name:           "GET Users",
			method:         "GET",
			path:           "/api/v1/user",
			expectedStatus: http.StatusOK,
			setupMocks: func(mockBooking *MockBookingUsecase, mockUser *MockUserUsecase, mockUnit *MockUnitUsecase) {
				mockUser.On("GetUsers", mock.Anything, 1, 10).Return([]domain.User{}, nil)
			},
		},
		{
			name:           "GET Units",
			method:         "GET",
			path:           "/api/v1/unit",
			expectedStatus: http.StatusOK,
			setupMocks: func(mockBooking *MockBookingUsecase, mockUser *MockUserUsecase, mockUnit *MockUnitUsecase) {
				mockUnit.On("GetUnits", mock.Anything, 1, 10).Return([]domain.Unit{}, nil)
			},
		},
		{
			name:           "GET Bookings",
			method:         "GET",
			path:           "/api/v1/booking",
			expectedStatus: http.StatusOK,
			setupMocks: func(mockBooking *MockBookingUsecase, mockUser *MockUserUsecase, mockUnit *MockUnitUsecase) {
				mockBooking.On("GetBookings", mock.Anything, 1, 10).Return([]domain.Booking{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBookingUsecase := new(MockBookingUsecase)
			mockUserUsecase := new(MockUserUsecase)
			mockUnitUsecase := new(MockUnitUsecase)
			handler := NewHandler(mockBookingUsecase, mockUserUsecase, mockUnitUsecase)

			tt.setupMocks(mockBookingUsecase, mockUserUsecase, mockUnitUsecase)

			req := httptest.NewRequest(tt.method, tt.path+"?page=1&page_size=10", nil)
			w := httptest.NewRecorder()

			switch tt.path {
			case "/api/v1/user":
				handler.GetUsers(w, req)
			case "/api/v1/unit":
				handler.GetUnits(w, req)
			case "/api/v1/booking":
				handler.GetBookings(w, req)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockBookingUsecase.AssertExpectations(t)
			mockUserUsecase.AssertExpectations(t)
			mockUnitUsecase.AssertExpectations(t)
		})
	}
}
