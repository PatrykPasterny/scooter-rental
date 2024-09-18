//go:build unit

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	mockredis "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/repository/mock"
	mockrental "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental/mock"
	rentalmodel "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental/model"
	mocktracker "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker/mock"
	trackermodel "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker/model"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/model"
)

const (
	testCity      = "Montreal"
	testLongitude = 70.0
	testLatitude  = 60.0
	testHeight    = 10000.0
	testWidth     = 15000.0
)

func TestGetScooters(t *testing.T) {
	s, mockRentalService, _ := beforeTest(t)

	ctx := context.Background()

	clientUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	rentalScooters := []*rentalmodel.Scooter{
		rentalmodel.NewScooter(
			scooterUUID.String(),
			testCity,
			testLongitude,
			testLatitude,
			true,
		),
	}

	expectedScooters := []model.ScooterGet{
		{
			ScooterUUID:  scooterUUID,
			Longitude:    testLongitude,
			Latitude:     testLatitude,
			Availability: true,
		},
	}

	params := &model.ScooterQueryParams{
		Longitude: testLongitude,
		Latitude:  testLatitude,
		Height:    testHeight,
		Width:     testWidth,
		City:      testCity,
	}

	rectangle := rentalmodel.NewRectangle(params.City, params.Longitude, params.Latitude, params.Height, params.Width)

	validURLQuery := &url.Values{}
	validURLQuery.Add("longitude", strconv.FormatFloat(params.Longitude, 'f', -1, 64))
	validURLQuery.Add("latitude", strconv.FormatFloat(params.Latitude, 'f', -1, 64))
	validURLQuery.Add("height", strconv.FormatFloat(params.Height, 'f', -1, 64))
	validURLQuery.Add("width", strconv.FormatFloat(params.Width, 'f', -1, 64))
	validURLQuery.Add("city", params.City)

	invalidURLQuery := &url.Values{}
	invalidURLQuery.Add("wrong", "wrong")

	expectedScootersJSON, err := json.Marshal(expectedScooters)
	require.NoError(t, err)

	tests := map[string]struct {
		mockRentalServiceHandler func(mock *mockrental.MockRentalService)
		urlQuery                 *url.Values
		clientUUID               uuid.NullUUID
		expectedCode             int
		expectedBody             string
	}{
		"successfully getting scooters": {
			mockRentalServiceHandler: func(mock *mockrental.MockRentalService) {
				mock.EXPECT().GetScooters(ctx, rectangle).
					Return(rentalScooters, nil).Times(1)
			},
			urlQuery:     validURLQuery,
			clientUUID:   uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode: http.StatusOK,
			expectedBody: string(expectedScootersJSON),
		},
		"failed getting scooter because request has no clientUUID in header": {
			mockRentalServiceHandler: nil,
			urlQuery:                 validURLQuery,
			clientUUID:               uuid.NullUUID{Valid: false},
			expectedCode:             http.StatusBadRequest,
			expectedBody:             `{"Message":"Failed getting clientUUID from header."}`,
		},
		"failed getting scooter because request has wrong query params": {
			mockRentalServiceHandler: nil,
			urlQuery:                 invalidURLQuery,
			clientUUID:               uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode:             http.StatusBadRequest,
			expectedBody:             `{"Message":"Failed decoding query params."}`,
		},
		"failed getting scooter because redis service threw error while getting scooters": {
			mockRentalServiceHandler: func(mock *mockrental.MockRentalService) {
				mock.EXPECT().GetScooters(ctx, rectangle).
					Return(nil, errors.New("")).Times(1)
			},
			urlQuery:     validURLQuery,
			clientUUID:   uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"Message":"Failed getting scooters."}`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request := buildRequest(t, scootersPath, http.MethodGet, &bytes.Buffer{}, tt.clientUUID)

			request.URL.RawQuery = tt.urlQuery.Encode()

			responseRecorder := httptest.NewRecorder()

			if tt.mockRentalServiceHandler != nil {
				tt.mockRentalServiceHandler(mockRentalService)
			}

			s.getScooters(responseRecorder, request)

			if status := responseRecorder.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got = %v want = %v",
					status, tt.expectedCode)
			}

			if body := responseRecorder.Body.String(); body != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got = %v want = %v",
					body, tt.expectedBody)
			}
		})
	}
}

func TestRentScooter(t *testing.T) {
	s, mockRentalService, mockTrackerService := beforeTest(t)

	ctx := context.Background()

	clientUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooter := model.RentPost{
		ScooterUUID: scooterUUID,
		Longitude:   testLongitude,
		Latitude:    testLatitude,
		City:        testCity,
	}

	scooterJSON, err := json.Marshal(scooter)
	require.NoError(t, err)

	invalidScooterJSON, err := json.Marshal("invalidScooter")
	require.NoError(t, err)

	rentInfo := rentalmodel.NewRentInfo(scooter.ScooterUUID.String(), scooter.City)
	trackerInfo := trackermodel.NewScooter(scooter.ScooterUUID.String(), scooter.City, scooter.Longitude, scooter.Latitude)

	tests := map[string]struct {
		mockRentalServiceHandler  func(mock *mockrental.MockRentalService)
		mockTrackerServiceHandler func(mock *mocktracker.MockService)
		body                      *bytes.Buffer
		clientUUID                uuid.NullUUID
		expectedCode              int
	}{
		"successfully renting scooter": {
			mockRentalServiceHandler: func(mock *mockrental.MockRentalService) {
				mock.EXPECT().Rent(ctx, rentInfo).Return(nil).Times(1)
			},
			mockTrackerServiceHandler: func(mock *mocktracker.MockService) {
				mock.EXPECT().Track(clientUUID, trackerInfo).Return(nil).Times(1)
			},
			body:         bytes.NewBuffer(scooterJSON),
			clientUUID:   uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode: http.StatusNoContent,
		},
		"failed renting scooter because request has no clientUUID in header": {
			mockRentalServiceHandler:  nil,
			mockTrackerServiceHandler: nil,
			body:                      bytes.NewBuffer(scooterJSON),
			clientUUID:                uuid.NullUUID{Valid: false},
			expectedCode:              http.StatusBadRequest,
		},
		"failed renting scooter because request has invalid body": {
			mockRentalServiceHandler:  nil,
			mockTrackerServiceHandler: nil,
			body:                      bytes.NewBuffer(invalidScooterJSON),
			clientUUID:                uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode:              http.StatusBadRequest,
		},
		"failed renting scooter because rental service threw error while renting scooter": {
			mockRentalServiceHandler: func(mock *mockrental.MockRentalService) {
				mock.EXPECT().Rent(ctx, rentInfo).Return(errors.New("")).Times(1)
			},
			mockTrackerServiceHandler: nil,
			body:                      bytes.NewBuffer(scooterJSON),
			clientUUID:                uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode:              http.StatusInternalServerError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request := buildRequest(t, rentPath, http.MethodPost, tt.body, tt.clientUUID)

			responseRecorder := httptest.NewRecorder()

			if tt.mockRentalServiceHandler != nil {
				tt.mockRentalServiceHandler(mockRentalService)
			}

			if tt.mockTrackerServiceHandler != nil {
				tt.mockTrackerServiceHandler(mockTrackerService)
			}

			s.rentScooter(responseRecorder, request)

			if status := responseRecorder.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got = %v want = %v",
					status, tt.expectedCode)
			}
		})
	}
}

func TestFreeScooter(t *testing.T) {
	s, mockRentalService, mockTrackerService := beforeTest(t)

	ctx := context.Background()

	clientUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooterUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	scooter := model.FreePost{
		ScooterUUID: scooterUUID,
	}

	scooterJSON, err := json.Marshal(scooter)
	require.NoError(t, err)

	invalidScooterJSON, err := json.Marshal("invalidScooter")
	require.NoError(t, err)

	tests := map[string]struct {
		mockRentalServiceHandler  func(mock *mockrental.MockRentalService)
		mockTrackerServiceHandler func(mock *mocktracker.MockService)
		body                      *bytes.Buffer
		clientUUID                uuid.NullUUID
		expectedCode              int
	}{
		"successfully freeing scooter": {
			mockRentalServiceHandler: func(mock *mockrental.MockRentalService) {
				mock.EXPECT().Free(ctx, scooterUUID).Return(nil).Times(1)
			},
			mockTrackerServiceHandler: func(mock *mocktracker.MockService) {
				mock.EXPECT().StopTracking(clientUUID, scooterUUID).Return(nil).Times(1)
			},
			body:         bytes.NewBuffer(scooterJSON),
			clientUUID:   uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode: http.StatusNoContent,
		},
		"failed freeing scooter because request has no clientUUID in header": {
			mockRentalServiceHandler:  nil,
			mockTrackerServiceHandler: nil,
			body:                      bytes.NewBuffer(scooterJSON),
			clientUUID:                uuid.NullUUID{Valid: false},
			expectedCode:              http.StatusBadRequest,
		},
		"failed freeing scooter because request has invalid body": {
			mockRentalServiceHandler:  nil,
			mockTrackerServiceHandler: nil,
			body:                      bytes.NewBuffer(invalidScooterJSON),
			clientUUID:                uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode:              http.StatusBadRequest,
		},
		"failed freeing scooter because rental service threw error while renting scooter": {
			mockRentalServiceHandler: func(mock *mockrental.MockRentalService) {
				mock.EXPECT().Free(ctx, scooterUUID).Return(errors.New("")).Times(1)
			},
			mockTrackerServiceHandler: nil,
			body:                      bytes.NewBuffer(scooterJSON),
			clientUUID:                uuid.NullUUID{UUID: clientUUID, Valid: true},
			expectedCode:              http.StatusInternalServerError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			request := buildRequest(t, freePath, http.MethodPost, tt.body, tt.clientUUID)

			responseRecorder := httptest.NewRecorder()

			if tt.mockRentalServiceHandler != nil {
				tt.mockRentalServiceHandler(mockRentalService)
			}

			if tt.mockTrackerServiceHandler != nil {
				tt.mockTrackerServiceHandler(mockTrackerService)
			}

			s.freeScooter(responseRecorder, request)

			if status := responseRecorder.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got = %v want = %v",
					status, tt.expectedCode)
			}
		})
	}
}

func beforeTest(t *testing.T) (*Server, *mockrental.MockRentalService, *mocktracker.MockService) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	validate := validator.New()

	controller := gomock.NewController(t)
	httpRouter := mux.NewRouter()

	mockRedisService := mockredis.NewMockScooterRepository(controller)
	mockRentalService := mockrental.NewMockRentalService(controller)
	mockTrackerService := mocktracker.NewMockService(controller)

	users := make(map[string]bool)

	s := NewServer(
		logger,
		validate,
		&http.Server{
			Addr:    fmt.Sprintf(":%d", 8081),
			Handler: httpRouter,
		},
		httpRouter,
		mockRedisService,
		mockRentalService,
		mockTrackerService,
		users,
	)

	return s, mockRentalService, mockTrackerService
}

func buildRequest(t *testing.T, path, method string, body *bytes.Buffer, clientUUID uuid.NullUUID) *http.Request {
	t.Helper()

	request, err := http.NewRequestWithContext(
		context.Background(),
		method,
		version+path,
		body,
	)
	require.NoErrorf(t, err, "Building new request")

	if clientUUID.Valid {
		request.Header.Set("Client-Id", clientUUID.UUID.String())
	}

	return request
}
