package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	rentalmodel "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const unitOfLength = "m" // in meters

var ErrScooterNotAvailable = errors.New("scooter with given ScooterUUID is not available")

func getScooters(
	ctx context.Context,
	client *redis.Client,
	geoRectangle *rentalmodel.GeoRectangle,
) ([]redis.GeoLocation, error) {
	// Perform the GeoRadius search
	scootersDB, err := client.GeoSearchLocation(ctx, geoRectangle.City, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude: geoRectangle.CenterLongitude,
			Latitude:  geoRectangle.CenterLatitude,
			BoxHeight: geoRectangle.Height,
			BoxWidth:  geoRectangle.Width,
			BoxUnit:   unitOfLength,
		},
		WithCoord: true,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("getting scooters from redis using geo search: %w", err)
	}

	return scootersDB, err
}

func getScooterAvailability(
	ctx context.Context,
	client *redis.Client,
	scooterUUID uuid.UUID,
) (bool, error) {
	// Retrieve the scooter directly using its ScooterUUID
	scooterJSON, err := client.Get(ctx, scooterUUID.String()).Result()
	if err != nil {
		return false, fmt.Errorf("getting scooters availability from redis: %w", err)
	}

	var availabilityAsInt int

	err = json.Unmarshal([]byte(scooterJSON), &availabilityAsInt)
	if err != nil {
		return false, fmt.Errorf("unmarshaling scooters availability: %w", err)
	}

	availability := availabilityAsInt == 1

	return availability, nil
}

func updateScooterLocation(
	ctx context.Context,
	client *redis.Client,
	scooter *redis.GeoLocation,
	city string,
) error {
	// Update the Geo index with scooter information
	if _, err := client.GeoAdd(ctx, city, scooter).Result(); err != nil {
		return fmt.Errorf("adding scooter's location to redis: %w", err)
	}

	return nil
}

func updateScooterAvailability(
	ctx context.Context,
	client *redis.Client,
	scooterUUID uuid.UUID,
	availability bool,
) error {
	key := scooterUUID.String()

	// make sure the availability is not changed from false to false or from true to true (from business side two
	// users won't be able to use the same scooter at the same time)
	if err := client.Watch(ctx, func(tx *redis.Tx) error {
		scooterAvailability, err := tx.Get(ctx, key).Result()
		if err != nil {
			return fmt.Errorf("getting scooter's availability from redis: %w", err)
		}

		if !willAvailabilityChange(availability, scooterAvailability) {
			return ErrScooterNotAvailable
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// Store the additional data as a string in Redis
			if err = pipe.Set(ctx, key, availability, 0).Err(); err != nil {
				return fmt.Errorf("updating scooter's availability in redis: %w", err)
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("error while executing the pipeline: %v", err)
		}

		return nil
	}, key); err != nil {
		return fmt.Errorf("updating scooter availability: %w", err)
	}

	return nil
}

func willAvailabilityChange(wantAvailability bool, redisAvailability string) bool {
	if (wantAvailability && redisAvailability == "1") || (!wantAvailability && redisAvailability == "0") {
		return false
	}

	return true
}
