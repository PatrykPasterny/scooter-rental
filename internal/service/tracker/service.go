package tracker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker/model"
)

const (
	north = 1
	west  = 2
	east  = 3
	south = 4

	MovingTimeInSeconds = 3

	oneSecondDecimal float64 = 0.000278
)

//go:generate mockgen -source=service.go -destination=mock/service_mock.go -package=mock
type Service interface {
	Track(userUUID uuid.UUID, scooter *model.Scooter) error
	StopTracking(userUUID, scooterUUID uuid.UUID) error
}

type trackingService struct {
	logger         *slog.Logger
	service        service.ScooterRepository
	rentedScooters map[uuid.UUID]chan uuid.UUID
	errorsChan     map[uuid.UUID]chan error
}

func NewTrackingService(logger *slog.Logger, service service.ScooterRepository) *trackingService {
	return &trackingService{
		logger:         logger,
		service:        service,
		rentedScooters: make(map[uuid.UUID]chan uuid.UUID),
		errorsChan:     make(map[uuid.UUID]chan error),
	}
}

// Track simulates the startup of a tracker go routine running on a scooter that periodically updates its localisation
// and also simulates its movement until the time the tracker go routine is stopped.
func (ts *trackingService) Track(userUUID uuid.UUID, scooter *model.Scooter) error {
	scooterUUID, err := uuid.Parse(scooter.Name)
	if err != nil {
		return fmt.Errorf("parsing scooter's uuid: %w", err)
	}

	trackerLogger := ts.logger.With(
		slog.String("scooter_id", scooterUUID.String()),
		slog.String("user_id", userUUID.String()),
	)

	currentScooterChan := make(chan uuid.UUID)
	currentErrorChan := make(chan error)

	if v, ok := ts.rentedScooters[scooterUUID]; ok && v != nil {
		close(v)
	}

	ts.rentedScooters[scooterUUID] = currentScooterChan
	ts.errorsChan[scooterUUID] = currentErrorChan

	trackerLogger.Info("Started tracking scooter")

	go func(tLogger *slog.Logger) {
		defer close(currentScooterChan)

		trackerContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		rentalErrors := make(map[string]int)
		for {
			select {
			case <-time.After(MovingTimeInSeconds * time.Second):
				simulateScooterMove(scooter, MovingTimeInSeconds, north)

				tLogger.Info(
					"Tracked scooter continues his journey.",
					slog.Float64("longitude", scooter.Longitude),
					slog.Float64("latitude", scooter.Latitude),
				)

				err = ts.service.UpdateScooterLocation(trackerContext, scooter)
				if err != nil {
					if _, ok := rentalErrors[err.Error()]; ok {
						rentalErrors[err.Error()] += 1
					}

					rentalErrors[err.Error()] = 1
				}
			case <-currentScooterChan: // Signal to stop tracking
				if len(rentalErrors) == 0 {
					currentErrorChan <- nil

					return
				}

				routineErrors := fmt.Errorf(
					"go routine assigned to userUUID - %s met several errors",
					scooterUUID,
				)

				for key, value := range rentalErrors {
					routineErrors = fmt.Errorf("%w: %s, %d times", routineErrors, key, value)
				}

				currentErrorChan <- routineErrors

				return
			}
		}
	}(trackerLogger)

	return nil
}

// StopTracking stops the tracking go routine for a given scooterUUID (simulates the stopping process on the scooter
// itself).
func (ts *trackingService) StopTracking(userUUID, scooterUUID uuid.UUID) error {
	defer close(ts.errorsChan[scooterUUID])

	ts.logger.Info(
		"Stopped tracking scooter.",
		slog.String("scooter_id", scooterUUID.String()),
		slog.String("user_id", userUUID.String()),
	)

	scooterToFree := ts.rentedScooters[scooterUUID]
	scooterToFree <- scooterUUID
	ts.rentedScooters[scooterUUID] = nil

	potentialErrors := <-ts.errorsChan[scooterUUID]
	if potentialErrors != nil {
		return fmt.Errorf("freeing scooter: %w", potentialErrors)
	}

	return nil
}

// simulateScooterMove is simulating the move of the scooter, I assume that each scooter goes on average 36 km/h
// which is around one second degree per second(approximately for both latitude and longitude). I pick
// one of four sides(north, west, east, south) and move the scooter three second degrees in that direction.
func simulateScooterMove(scooter *model.Scooter, timeInSeconds int, direction int) {
	if direction == north {
		scooter.Latitude += float64(timeInSeconds) * oneSecondDecimal
	} else if direction == south {
		scooter.Latitude -= float64(timeInSeconds) * oneSecondDecimal
	} else if direction == east {
		scooter.Longitude += float64(timeInSeconds) * oneSecondDecimal
	} else if direction == west {
		scooter.Longitude -= float64(timeInSeconds) * oneSecondDecimal
	}
}
