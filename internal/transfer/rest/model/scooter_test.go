//go:build unit

package model

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	testLongitude = 70.0
	testLatitude  = 60.0
)

func TestFilterScooters(t *testing.T) {
	firstScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	secondScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	thirdScooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scootersToFilter := []ScooterGet{
		{
			ScooterUUID:  firstScooterUUID,
			Longitude:    testLongitude,
			Latitude:     testLatitude,
			Availability: true,
		},
		{
			ScooterUUID:  secondScooterUUID,
			Longitude:    testLongitude,
			Latitude:     testLatitude,
			Availability: false,
		},
		{
			ScooterUUID:  thirdScooterUUID,
			Longitude:    testLongitude,
			Latitude:     testLatitude,
			Availability: false,
		},
	}

	availableScooters := []ScooterGet{
		scootersToFilter[0],
	}

	nonAvailableScooters := []ScooterGet{
		scootersToFilter[1],
		scootersToFilter[2],
	}

	tests := map[string]struct {
		f    func(s *ScooterGet) bool
		want []ScooterGet
	}{
		"Successfully returned available scooters": {
			f: func(s *ScooterGet) bool {
				return s.Availability
			},
			want: availableScooters,
		},
		"Successfully returned non available scooters": {
			f: func(s *ScooterGet) bool {
				return !s.Availability
			},
			want: nonAvailableScooters,
		},
	}
	for tName, tt := range tests {
		t.Run(tName, func(t *testing.T) {
			if got := FilterScooters(scootersToFilter, tt.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterScooters() = %v, want %v", got, tt.want)
			}
		})
	}
}
