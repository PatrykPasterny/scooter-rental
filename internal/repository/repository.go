package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	rentalmodel "github.com/PatrykPasterny/scooter-rental/internal/service/rental/model"
	trackermodel "github.com/PatrykPasterny/scooter-rental/internal/service/tracker/model"
)

type redisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *redisService {
	return &redisService{
		client: client,
	}
}

func (rs *redisService) GetScooters(
	ctx context.Context,
	geoRectangle *rentalmodel.GeoRectangle,
) ([]*rentalmodel.Scooter, error) {
	var results []*rentalmodel.Scooter

	scooters, err := getScooters(ctx, rs.client, geoRectangle)
	if err != nil {
		return nil, fmt.Errorf("getting scooters: %w", err)
	}

	results = make([]*rentalmodel.Scooter, len(scooters))

	for i := range scooters {
		var scooterUUID uuid.UUID

		scooterUUID, innerErr := uuid.Parse(scooters[i].Name)
		if innerErr != nil {
			return nil, fmt.Errorf("parsing scooter's uuid: %w", innerErr)
		}

		var availability bool

		availability, innerErr = getScooterAvailability(ctx, rs.client, scooterUUID)
		if innerErr != nil {
			return nil, fmt.Errorf("getting scooter's availability: %w", innerErr)
		}

		result := rentalmodel.NewScooter(
			scooters[i].Name,
			geoRectangle.City,
			scooters[i].Longitude,
			scooters[i].Latitude,
			availability,
		)
		results[i] = result
	}

	return results, nil
}

func (rs *redisService) UpdateScooterLocation(ctx context.Context, scooter *trackermodel.Scooter) error {
	redisLocation := &redis.GeoLocation{
		Name:      scooter.Name,
		Longitude: scooter.Longitude,
		Latitude:  scooter.Latitude,
	}

	err := updateScooterLocation(ctx, rs.client, redisLocation, scooter.City)
	if err != nil {
		return fmt.Errorf("updating scooter's location: %w", err)
	}

	return nil
}

func (rs *redisService) UpdateScooterAvailability(ctx context.Context, scooterUUID uuid.UUID, availability bool) error {
	err := updateScooterAvailability(ctx, rs.client, scooterUUID, availability)
	if err != nil {
		return fmt.Errorf("updating scooter's availability: %w", err)
	}

	return nil
}
