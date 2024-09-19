//go:build unit

package rental

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	repositorymock "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/mock"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental/model"
)

const (
	testCity = "Montreal"
)

func TestRent(t *testing.T) {
	ctx := context.Background()

	firstScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	rentInfo := model.NewRentInfo(
		firstScooterUUID.String(),
		testCity,
	)

	wrongUUIDScooter := model.NewRentInfo(
		"dd-dd-dd",
		testCity,
	)

	tests := map[string]struct {
		rentInfo                *model.RentInfo
		mockRedisServiceHandler func(mock *repositorymock.MockScooterRepository)
		wantErr                 bool
	}{
		"successfully rent scooter": {
			rentInfo: rentInfo,
			mockRedisServiceHandler: func(mock *repositorymock.MockScooterRepository) {
				scooterUUID, innerErr := uuid.Parse(rentInfo.ScooterUUID)
				require.NoError(t, innerErr)

				mock.EXPECT().UpdateScooterAvailability(ctx, scooterUUID, false).Return(nil).Times(1)
			},
			wantErr: false,
		},
		"rent scooter failing because chosen scooter has incorrect ScooterUUID": {
			rentInfo:                wrongUUIDScooter,
			mockRedisServiceHandler: nil,
			wantErr:                 true,
		},
		"rent scooter failing because redis service threw an error": {
			rentInfo: rentInfo,
			mockRedisServiceHandler: func(mock *repositorymock.MockScooterRepository) {
				scooterUUID, innerErr := uuid.Parse(rentInfo.ScooterUUID)
				require.NoError(t, innerErr)

				mock.EXPECT().UpdateScooterAvailability(ctx, scooterUUID, false).Return(redis.ErrClosed)
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRedisService := repositorymock.NewMockScooterRepository(controller)

			if tt.mockRedisServiceHandler != nil {
				tt.mockRedisServiceHandler(mockRedisService)
			}

			rs := NewRentalService(mockRedisService)
			if err = rs.Rent(ctx, tt.rentInfo); (err != nil) != tt.wantErr {
				t.Errorf("Rent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFree(t *testing.T) {
	ctx := context.Background()

	logger := log.New(os.Stdout, "TEST ", log.LstdFlags)

	firstScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	tests := map[string]struct {
		logger                  *log.Logger
		mockRedisServiceHandler func(mock *repositorymock.MockScooterRepository)
		wantErr                 bool
	}{
		"successfully freed scooter": {
			logger: logger,
			mockRedisServiceHandler: func(mock *repositorymock.MockScooterRepository) {
				mock.EXPECT().UpdateScooterAvailability(ctx, firstScooterUUID, true).
					Return(nil).Times(1)
			},
			wantErr: false,
		},
		"freeing scooter failed because redis service threw an error when updating availability": {
			logger: logger,
			mockRedisServiceHandler: func(mock *repositorymock.MockScooterRepository) {
				mock.EXPECT().UpdateScooterAvailability(ctx, firstScooterUUID, true).
					Return(redis.ErrClosed).Times(1)
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRedisService := repositorymock.NewMockScooterRepository(controller)

			if tt.mockRedisServiceHandler != nil {
				tt.mockRedisServiceHandler(mockRedisService)
			}

			rs := NewRentalService(mockRedisService)
			if err = rs.Free(ctx, firstScooterUUID); (err != nil) != tt.wantErr {
				t.Errorf("Free() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
