//go:build unit

package repository

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	rentalmodel "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental/model"
	trackermodel "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker/model"
)

const (
	testCity      = "Montreal"
	testLongitude = 70.0
	testLatitude  = 60.0
	testHeight    = 10000.0
	testWidth     = 15000.0
)

func TestGetScootersRepo(t *testing.T) {
	ctx := context.Background()

	firstScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	secScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	thirdScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scootersInRectangle := []redis.GeoLocation{
		{
			Name:      firstScooterUUID.String(),
			Longitude: 60.0,
			Latitude:  40.0,
		},
		{
			Name:      secScooterUUID.String(),
			Longitude: 60.001,
			Latitude:  40.02,
		},
		{
			Name:      thirdScooterUUID.String(),
			Longitude: 60.12415,
			Latitude:  -40.0,
		},
	}

	scootersActivities := []string{"1", "0", "1"}
	scootersActivitiesAsBool := []bool{true, false, true}

	scooters := []*rentalmodel.Scooter{
		rentalmodel.NewScooter(
			scootersInRectangle[0].Name,
			testCity,
			scootersInRectangle[0].Longitude,
			scootersInRectangle[0].Latitude,
			scootersActivitiesAsBool[0],
		),
		rentalmodel.NewScooter(
			scootersInRectangle[1].Name,
			testCity,
			scootersInRectangle[1].Longitude,
			scootersInRectangle[1].Latitude,
			scootersActivitiesAsBool[1],
		),
		rentalmodel.NewScooter(
			scootersInRectangle[2].Name,
			testCity,
			scootersInRectangle[2].Longitude,
			scootersInRectangle[2].Latitude,
			scootersActivitiesAsBool[2],
		),
	}

	geoRectangle := rentalmodel.NewRectangle(testCity, testLongitude, testLatitude, testHeight, testWidth)

	tests := map[string]struct {
		redisMock func(mock redismock.ClientMock)
		want      []*rentalmodel.Scooter
		wantErr   bool
	}{
		"getting scooters successfully": {
			redisMock: func(mock redismock.ClientMock) {
				mock.ExpectGeoSearchLocation(testCity, &redis.GeoSearchLocationQuery{
					GeoSearchQuery: redis.GeoSearchQuery{
						Longitude: testLongitude,
						Latitude:  testLatitude,
						BoxHeight: testHeight,
						BoxWidth:  testWidth,
						BoxUnit:   unitOfLength,
					},
					WithCoord: true,
				}).SetVal(scootersInRectangle)

				for i := range scootersInRectangle {
					mock.ExpectGet(scootersInRectangle[i].Name).SetVal(scootersActivities[i])
				}
			},
			want:    scooters,
			wantErr: false,
		},
		"getting scooters failed, because repository threw an error when getting scooters": {
			redisMock: func(mock redismock.ClientMock) {
				mock.ExpectGeoSearchLocation(testCity, &redis.GeoSearchLocationQuery{
					GeoSearchQuery: redis.GeoSearchQuery{
						Longitude: testLongitude,
						Latitude:  testLatitude,
						BoxHeight: testHeight,
						BoxWidth:  testWidth,
						BoxUnit:   unitOfLength,
					},
					WithCoord: true,
				}).SetErr(redis.ErrClosed)
			},
			want:    nil,
			wantErr: true,
		},
		"getting scooters failed, because repository threw an error when getting scooter's availability": {
			redisMock: func(mock redismock.ClientMock) {
				mock.ExpectGeoSearchLocation(testCity, &redis.GeoSearchLocationQuery{
					GeoSearchQuery: redis.GeoSearchQuery{
						Longitude: testLongitude,
						Latitude:  testLatitude,
						BoxHeight: testHeight,
						BoxWidth:  testWidth,
						BoxUnit:   unitOfLength,
					},
					WithCoord: true,
				}).SetVal(scootersInRectangle)

				mock.ExpectGet(scootersInRectangle[0].Name).SetErr(redis.ErrClosed)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			redisClient, redisMock := redismock.NewClientMock()

			tt.redisMock(redisMock)

			rs := NewRedisService(redisClient)

			got, err := rs.GetScooters(ctx, geoRectangle)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScooters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetScooters() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateScooterLocation(t *testing.T) {
	ctx := context.Background()

	logger := &log.Logger{}

	scooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooter := &redis.GeoLocation{
		Name:      scooterUUID.String(),
		Longitude: testLongitude,
		Latitude:  testLatitude,
	}

	trackerScooter := trackermodel.NewScooter(scooter.Name, testCity, scooter.Longitude, scooter.Latitude)

	tests := map[string]struct {
		logger                   *log.Logger
		mockRedisDatabaseHandler func(mock redismock.ClientMock)
		wantErr                  bool
	}{
		"updating scooter successfully": {
			logger: logger,
			mockRedisDatabaseHandler: func(mock redismock.ClientMock) {
				mock.ExpectGeoAdd(testCity, scooter).SetVal(1)
			},
			wantErr: false,
		},
		"updating scooter failed, because repository threw an error when updating scooter's location": {
			logger: logger,
			mockRedisDatabaseHandler: func(mock redismock.ClientMock) {
				mock.ExpectGeoAdd(testCity, scooter).SetErr(redis.ErrClosed)
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			redisClient, redisMock := redismock.NewClientMock()

			tt.mockRedisDatabaseHandler(redisMock)

			rs := NewRedisService(redisClient)

			if err = rs.UpdateScooterLocation(ctx, trackerScooter); (err != nil) != tt.wantErr {
				t.Errorf("UpdateScooterLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateScooterAvailability(t *testing.T) {
	ctx := context.Background()

	logger := &log.Logger{}

	firstScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooterAvailability := true

	tests := map[string]struct {
		logger                   *log.Logger
		mockRedisDatabaseHandler func(mock redismock.ClientMock)
		wantErr                  bool
	}{
		"updating scooter availability successfully": {
			logger: logger,
			mockRedisDatabaseHandler: func(mock redismock.ClientMock) {
				mock.ExpectWatch(firstScooterUUID.String())
				mock.ExpectGet(firstScooterUUID.String()).SetVal("0")
				mock.ExpectTxPipeline()
				mock.ExpectSet(firstScooterUUID.String(), scooterAvailability, 0).SetVal("status")
				mock.ExpectTxPipelineExec()
			},
			wantErr: false,
		},
		"updating scooter failed, because scooter was already available": {
			logger: logger,
			mockRedisDatabaseHandler: func(mock redismock.ClientMock) {
				mock.ExpectWatch(firstScooterUUID.String()).SetErr(ErrScooterNotAvailable)
				mock.ExpectGet(firstScooterUUID.String()).SetVal("1")
				mock.ExpectSet(firstScooterUUID.String(), scooterAvailability, 0)
				mock.ExpectTxPipelineExec()
			},
			wantErr: true,
		},
		"updating scooter failed, because repository threw an error when getting scooter's availability": {
			logger: logger,
			mockRedisDatabaseHandler: func(mock redismock.ClientMock) {
				mock.ExpectWatch(firstScooterUUID.String())
				mock.ExpectGet(firstScooterUUID.String()).SetErr(redis.ErrClosed)
			},
			wantErr: true,
		},
		"updating scooter failed, because repository threw an error when updating scooter's availability": {
			logger: logger,
			mockRedisDatabaseHandler: func(mock redismock.ClientMock) {
				mock.ExpectWatch(firstScooterUUID.String())
				mock.ExpectGet(firstScooterUUID.String())
				mock.ExpectSet(firstScooterUUID.String(), scooterAvailability, 0).SetErr(redis.ErrClosed)
				mock.ExpectTxPipelineExec()
			},
			wantErr: true,
		},
		"updating scooter failed, because repository threw an error when executing redis commands in pipeline": {
			logger: logger,
			mockRedisDatabaseHandler: func(mock redismock.ClientMock) {
				mock.ExpectTxPipeline()
				mock.ExpectGet(firstScooterUUID.String())
				mock.ExpectSet(firstScooterUUID.String(), scooterAvailability, 0).SetVal("status")
				mock.ExpectTxPipelineExec().SetErr(redis.ErrClosed)
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			redisClient, redisMock := redismock.NewClientMock()

			tt.mockRedisDatabaseHandler(redisMock)

			rs := NewRedisService(redisClient)
			if err = rs.UpdateScooterAvailability(ctx, firstScooterUUID, scooterAvailability); (err != nil) != tt.wantErr {
				t.Errorf("UpdateScooterAvailability() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
