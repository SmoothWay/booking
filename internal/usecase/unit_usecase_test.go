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

type MockUnitRepository struct {
	mock.Mock
}

func (m *MockUnitRepository) GetUnitByID(ctx context.Context, id string) (*domain.Unit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Unit), args.Error(1)
}

func (m *MockUnitRepository) CreateUnit(ctx context.Context, unit *domain.Unit) error {
	args := m.Called(ctx, unit)
	return args.Error(0)
}

func (m *MockUnitRepository) UpdateUnit(ctx context.Context, unit *domain.Unit) error {
	args := m.Called(ctx, unit)
	return args.Error(0)
}

func (m *MockUnitRepository) DeleteUnit(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUnitRepository) GetUnits(ctx context.Context, page int, pageSize int) ([]*domain.Unit, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Unit), args.Error(1)
}

func TestUnitUsecase_GetUnitByID_Success(t *testing.T) {
	mockRepo := new(MockUnitRepository)
	usecase := NewUnitUsecase(nil, mockRepo)

	expectedUnit := &domain.Unit{
		ID:          "1",
		Name:        "Test Unit",
		Description: "Test Description",
		Price:       100,
		CreatedAt:   time.Now(),
	}

	mockRepo.On("GetUnitByID", mock.Anything, "1").Return(expectedUnit, nil)

	result, err := usecase.GetUnitByID(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, expectedUnit.ID, result.ID)
	assert.Equal(t, expectedUnit.Name, result.Name)
	assert.Equal(t, expectedUnit.Description, result.Description)
	assert.Equal(t, expectedUnit.Price, result.Price)

	mockRepo.AssertExpectations(t)
}

func TestUnitUsecase_GetUnitByID_NotFound(t *testing.T) {
	mockRepo := new(MockUnitRepository)
	usecase := NewUnitUsecase(nil, mockRepo)

	expectedError := errors.New("unit not found")
	mockRepo.On("GetUnitByID", mock.Anything, "999").Return(nil, expectedError)

	result, err := usecase.GetUnitByID(context.Background(), "999")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, domain.Unit{}, result)

	mockRepo.AssertExpectations(t)
}

func TestUnitUsecase_GetUnits_Success(t *testing.T) {
	mockRepo := new(MockUnitRepository)
	usecase := NewUnitUsecase(nil, mockRepo)

	expectedUnits := []*domain.Unit{
		{
			ID:          "1",
			Name:        "Unit 1",
			Description: "Description 1",
			Price:       100,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "Unit 2",
			Description: "Description 2",
			Price:       200,
			CreatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetUnits", mock.Anything, 1, 10).Return(expectedUnits, nil)

	result, err := usecase.GetUnits(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedUnits[0].ID, result[0].ID)
	assert.Equal(t, expectedUnits[1].ID, result[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestUnitUsecase_GetUnits_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockUnitRepository)
	usecase := NewUnitUsecase(nil, mockRepo)

	expectedError := errors.New("database error")
	mockRepo.On("GetUnits", mock.Anything, 1, 10).Return(nil, expectedError)

	// Act
	result, err := usecase.GetUnits(context.Background(), 1, 10)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestUnitUsecase_CreateUnit_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUnitRepository)
	usecase := NewUnitUsecase(nil, mockRepo)

	event := events.UnitCreatedEvent{
		ID:          "1",
		Name:        "New Unit",
		Description: "New Description",
		Price:       150,
		CreatedAt:   time.Now(),
	}

	mockRepo.On("CreateUnit", mock.Anything, mock.MatchedBy(func(unit *domain.Unit) bool {
		return unit.ID == event.ID && unit.Name == event.Name
	})).Return(nil)

	// Act
	err := usecase.CreateUnit(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUnitUsecase_CreateUnit_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockUnitRepository)
	usecase := NewUnitUsecase(nil, mockRepo)

	event := events.UnitCreatedEvent{
		ID:          "1",
		Name:        "New Unit",
		Description: "New Description",
		Price:       150,
		CreatedAt:   time.Now(),
	}

	expectedError := errors.New("database error")
	mockRepo.On("CreateUnit", mock.Anything, mock.Anything).Return(expectedError)

	// Act
	err := usecase.CreateUnit(context.Background(), event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestUnitUsecase_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		unitID      string
		expectError bool
		setupMock   func(*MockUnitRepository)
	}{
		{
			name:        "Empty ID",
			unitID:      "",
			expectError: true,
			setupMock: func(mockRepo *MockUnitRepository) {
				mockRepo.On("GetUnitByID", mock.Anything, "").Return(nil, errors.New("invalid id"))
			},
		},
		{
			name:        "Very Long ID",
			unitID:      "very-long-id-that-exceeds-normal-length-limits",
			expectError: false,
			setupMock: func(mockRepo *MockUnitRepository) {
				unit := &domain.Unit{ID: "very-long-id-that-exceeds-normal-length-limits"}
				mockRepo.On("GetUnitByID", mock.Anything, "very-long-id-that-exceeds-normal-length-limits").Return(unit, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockUnitRepository)
			usecase := NewUnitUsecase(nil, mockRepo)

			tt.setupMock(mockRepo)

			// Act
			_, err := usecase.GetUnitByID(context.Background(), tt.unitID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
