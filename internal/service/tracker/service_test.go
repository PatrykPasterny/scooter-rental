//go:build unit

package tracker

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/repository/mock"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker/model"
)

const (
	amountOfScooterTrackingEvents = 2
	firstTestCity                 = "Montreal"
	secondTestCity                = "Ottawa"
)

func TestTrackScooter(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	userUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	firstScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	thirdScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooters := []*model.Scooter{
		{
			Name:      firstScooterUUID.String(),
			Longitude: 70.01,
			Latitude:  60.01,
			City:      firstTestCity,
		},
		{
			Name:      thirdScooterUUID.String(),
			Longitude: 69.99,
			Latitude:  59.99,
			City:      secondTestCity,
		},
	}

	tests := map[string]struct {
		logger                  *slog.Logger
		mockRedisServiceHandler func(mock *mock.MockScooterRepository)
		wantErr                 bool
	}{
		"successfully tracking multiple scooters": {
			logger: logger,
			mockRedisServiceHandler: func(mock *mock.MockScooterRepository) {
				for i := range scooters {
					mock.EXPECT().UpdateScooterLocation(gomock.Any(), scooters[i]).
						Return(nil).Times(amountOfScooterTrackingEvents)
				}
			},
			wantErr: false,
		},
		"failed tracking multiple scooters, because of redis service threw error when updating scooter location ": {
			logger: logger,
			mockRedisServiceHandler: func(mock *mock.MockScooterRepository) {
				mock.EXPECT().UpdateScooterLocation(gomock.Any(), scooters[0]).
					Return(nil).Times(amountOfScooterTrackingEvents - 1)
				mock.EXPECT().UpdateScooterLocation(gomock.Any(), scooters[0]).
					Return(redis.ErrClosed).Times(1)
				for i := 1; i < len(scooters); i++ {
					mock.EXPECT().UpdateScooterLocation(gomock.Any(), scooters[i]).
						Return(nil).Times(amountOfScooterTrackingEvents)
				}
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRedisService := mock.NewMockScooterRepository(controller)

			tt.mockRedisServiceHandler(mockRedisService)

			ts := NewTrackingService(tt.logger, mockRedisService)

			for i := range scooters {
				innerErr := ts.Track(userUUID, scooters[i])
				require.NoError(t, innerErr)
			}

			if len(ts.rentedScooters) != len(scooters) {
				t.Errorf(
					"Track() should rent all given scooters = %v, want %v",
					len(ts.rentedScooters),
					len(scooters),
				)
			}

			for i := 0; i < amountOfScooterTrackingEvents; i++ {
				time.Sleep((MovingTimeInSeconds + 1) * time.Second)
			}

			for i := range scooters {
				scooterUUID, innerErr := uuid.Parse(scooters[i].Name)
				require.NoError(t, innerErr)

				ts.rentedScooters[scooterUUID] <- scooterUUID
			}

			var errorFound bool

			for i := range ts.errorsChan {
				if err = <-ts.errorsChan[i]; err != nil {
					errorFound = true
				}

				close(ts.errorsChan[i])
			}

			if errorFound != tt.wantErr {
				t.Errorf("Track() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func TestFreeScooter(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	userUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	firstScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooter := &model.Scooter{
		Name:      firstScooterUUID.String(),
		Longitude: 70.01,
		Latitude:  60.01,
		City:      firstTestCity,
	}

	tests := map[string]struct {
		logger                  *slog.Logger
		mockRedisServiceHandler func(mock *mock.MockScooterRepository)
		rentScooterHandler      func(tracker *trackingService) error
		wantErr                 bool
	}{
		"successfully freeing scooter": {
			logger:                  logger,
			mockRedisServiceHandler: nil,
			rentScooterHandler: func(ts *trackingService) error {
				return ts.Track(firstScooterUUID, scooter)
			},
			wantErr: false,
		},
		"freeing scooter failed, because scooter's rental process threw error": {
			logger: logger,
			mockRedisServiceHandler: func(mock *mock.MockScooterRepository) {
				mock.EXPECT().UpdateScooterLocation(gomock.Any(), scooter).Return(redis.ErrClosed)
			},
			rentScooterHandler: func(ts *trackingService) error {
				innerErr := ts.Track(firstScooterUUID, scooter)
				require.NoError(t, innerErr)

				time.Sleep((MovingTimeInSeconds + 1) * time.Second)

				return nil
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRedisService := mock.NewMockScooterRepository(controller)
			if tt.mockRedisServiceHandler != nil {
				tt.mockRedisServiceHandler(mockRedisService)
			}

			ts := NewTrackingService(tt.logger, mockRedisService)

			err = tt.rentScooterHandler(ts)
			require.NoError(t, err)

			if err = ts.StopTracking(userUUID, firstScooterUUID); (err != nil) != tt.wantErr {
				t.Errorf("StopTracking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
