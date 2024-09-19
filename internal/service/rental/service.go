package rental

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental/model"
)

//go:generate mockgen -source=service.go -destination=mock/service_mock.go -package=mock
type RentalService interface {
	GetScooters(ctx context.Context, rectangle *model.GeoRectangle) ([]*model.Scooter, error)
	Rent(ctx context.Context, info *model.RentInfo) error
	Free(ctx context.Context, scooterUUID uuid.UUID) error
}

type rentalService struct {
	scooterRepository service.ScooterRepository
}

func NewRentalService(repo service.ScooterRepository) *rentalService {
	return &rentalService{
		scooterRepository: repo,
	}
}

func (rs *rentalService) GetScooters(ctx context.Context, rectangle *model.GeoRectangle) ([]*model.Scooter, error) {
	scooters, err := rs.scooterRepository.GetScooters(ctx, rectangle)
	if err != nil {
		return nil, fmt.Errorf("getting scooters in the searched area: %w", err)
	}

	return scooters, err
}

func (rs *rentalService) Rent(ctx context.Context, info *model.RentInfo) error {
	scooterUUID, err := uuid.Parse(info.ScooterUUID)
	if err != nil {
		return fmt.Errorf("parsing scooter's uuid: %w", err)
	}

	err = rs.scooterRepository.UpdateScooterAvailability(ctx, scooterUUID, false)
	if err != nil {
		return fmt.Errorf("updating scooter availability: %w", err)
	}

	return nil
}

func (rs *rentalService) Free(ctx context.Context, scooterUUID uuid.UUID) error {
	if err := rs.scooterRepository.UpdateScooterAvailability(ctx, scooterUUID, true); err != nil {
		return fmt.Errorf("updating scooter availability: %w", err)
	}

	return nil
}
