package service

import (
	"context"

	"github.com/google/uuid"

	rentalmodel "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental/model"
	trackermodel "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker/model"
)

//go:generate mockgen -source=scooter_repository.go -destination=mock/scooter_repository_mock.go -package=mock
type ScooterRepository interface {
	GetScooters(ctx context.Context, geoRectangle *rentalmodel.GeoRectangle) ([]*rentalmodel.Scooter, error)
	UpdateScooterLocation(ctx context.Context, scooter *trackermodel.Scooter) error
	UpdateScooterAvailability(ctx context.Context, scooterUUID uuid.UUID, availability bool) error
}
