package usecase

import (
	"context"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/broker/kafka"
)

type UnitUsecase struct {
	kafkaProducer *kafka.Producer
	unitRepo      domain.UnitRepository
}

func NewUnitUsecase(kafkaProducer *kafka.Producer, unitRepo domain.UnitRepository) *UnitUsecase {
	return &UnitUsecase{
		kafkaProducer: kafkaProducer,
		unitRepo:      unitRepo,
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

func (uc *UnitUsecase) GetUnits(ctx context.Context, page int, pageSize int) ([]domain.Unit, error) {
	units, err := uc.unitRepo.GetUnits(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var unitsDomain []domain.Unit
	for _, unit := range units {
		unitsDomain = append(unitsDomain, *unit)
	}
	return unitsDomain, nil
}

func (uc *UnitUsecase) GetUnitByID(ctx context.Context, id string) (domain.Unit, error) {
	unit, err := uc.unitRepo.GetUnitByID(ctx, id)
	if err != nil {
		return domain.Unit{}, err
	}
	return *unit, nil
}
