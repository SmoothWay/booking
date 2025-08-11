package usecase

import (
	"context"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
)

type UnitUsecase struct {
	unitRepo domain.UnitRepository
}

func NewUnitUsecase(unitRepo domain.UnitRepository) *UnitUsecase {
	return &UnitUsecase{
		unitRepo: unitRepo,
	}
}

// CreateUnit handles unit creation from Kafka events
func (uc *UnitUsecase) CreateUnit(ctx context.Context, event events.UnitCreatedEvent) error {
	unit := &domain.Unit{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Price:       event.Price,
		CreatedAt:   event.CreatedAt,
	}

	return uc.unitRepo.CreateUnit(ctx, unit)
}
