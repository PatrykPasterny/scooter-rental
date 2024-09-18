package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/schema"

	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/model"
	modelrental "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental/model"
	trackermodel "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker/model"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

var (
	errExpectedHeaderParamNotFound = errors.New("expected header parameter was not found")
)

//	@title			Scootin Aboot
//	@version		1.0
//	@description	The API that enables user to rent, track and free scooters operated by Scootin Aboot company.
//	@Schema			http https
//	@BasePath		/api/v1

// getScooters returns all the scooters owned by Scootin Aboot company in the queried rectangle area of a given city.
//
//	@Summary	Gets scooters in the queried area of given city.
//	@Tags		scooters
//
//	@Param		Client-Id	header		string	true	"ClientID"									minlength(36)	maxlength(36)	default(00000000-0000-0000-0000-000000000000)
//	@Param		city		query		string	true	"City"										default(Ottawa)
//	@Param		longitude	query		number	true	"Longitude of the center of the rectangle"	default(73.4)
//	@Param		latitude	query		number	true	"Latitude of the center of the rectangle"	default(45.4)
//	@Param		height		query		number	true	"Height of the rectangle in meters"			default(20000.0)
//	@Param		width		query		number	true	"Width of the rectangle in meters"			default(25000.0)
//
//	@Success	200			{object}	[]model.ScooterGet
//	@Failure	400			{object}	model.ApiError
//	@Failure	403			{object}	model.ApiError
//	@Failure	500			{object}	model.ApiError
//	@Router		/scooters [get]
func (s *Server) getScooters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clientUUID, err := clientUUIDFromHeader(r)
	if err != nil {
		s.logger.Error("failed to get clientID from header", err)

		Error(w, http.StatusBadRequest, "Failed getting clientUUID from header.")

		return
	}

	ctxLogger := s.logger.With(
		slog.String("client_id", clientUUID.String()),
	)

	var queryParams model.ScooterQueryParams

	decoder := schema.NewDecoder()

	if err = decoder.Decode(&queryParams, r.URL.Query()); err != nil {
		ctxLogger.Error("failed to decode query params", err)

		Error(w, http.StatusBadRequest, "Failed decoding query params.")

		return
	}

	if err = s.validator.Struct(queryParams); err != nil {
		ctxLogger.Error("failed to validate query params", err)

		Error(w, http.StatusBadRequest, "Failed validating query params.")

		return
	}

	geoRectangle := modelrental.NewRectangle(
		queryParams.City,
		queryParams.Longitude,
		queryParams.Latitude,
		queryParams.Height,
		queryParams.Width,
	)

	ctxLogger = s.logger.With(
		slog.String("city", queryParams.City),
		slog.Float64("longitude", queryParams.Longitude),
		slog.Float64("latitude", queryParams.Latitude),
		slog.Float64("height", queryParams.Height),
		slog.Float64("width", queryParams.Width),
	)

	ctxLogger.Info("getting scooters")

	rentalScooters, err := s.rentalService.GetScooters(ctx, geoRectangle)
	if err != nil {
		ctxLogger.Error("failed to get scooters from rental service", err)

		Error(w, http.StatusInternalServerError, "Failed getting scooters.")

		return
	}

	scooters := make([]model.ScooterGet, len(rentalScooters))

	for i := range rentalScooters {
		scooterUUID, innerErr := uuid.Parse(rentalScooters[i].Name)
		if innerErr != nil {
			ctxLogger.Error("failed to get parse scooterID to ScooterUUID", err)

			Error(w, http.StatusInternalServerError, "Failed parsing scooterID.")

			return
		}

		scooters[i] = model.ScooterGet{
			ScooterUUID:  scooterUUID,
			Latitude:     rentalScooters[i].Latitude,
			Longitude:    rentalScooters[i].Longitude,
			Availability: rentalScooters[i].Availability,
		}
	}

	ctxLogger.Info("successfully received scooters")

	JSON(w, http.StatusOK, scooters)
}

// rentScooter enables user to rent the given scooter from the pool owned by Scootin Aboot company in a given city.
//
//	@Summary	Rents the chosen scooter in given city.
//	@Tags		scooters
//
//	@Param		Client-Id	header	string			true	"ClientID"	minlength(36)	maxlength(36)	default(00000000-0000-0000-0000-000000000000)
//	@Param		Payload		body	model.RentPost	true	"Rental information details"
//
//	@Success	204
//	@Failure	400	{object}	model.ApiError
//	@Failure	403	{object}	model.ApiError
//	@Failure	500	{object}	model.ApiError
//	@Router		/rent [post]
func (s *Server) rentScooter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clientUUID, err := clientUUIDFromHeader(r)
	if err != nil {
		s.logger.Error("failed to get clientID from header", err)

		Error(w, http.StatusBadRequest, "Failed getting clientUUID from header.")

		return
	}

	ctxLogger := s.logger.With(
		slog.String("client_id", clientUUID.String()),
	)

	var rentPost model.RentPost

	if err = json.NewDecoder(r.Body).Decode(&rentPost); err != nil {
		ctxLogger.Error("failed to decode request body", err)

		Error(w, http.StatusBadRequest, "Failed to decode request body to rental information.")

		return
	}

	if err = s.validator.Struct(rentPost); err != nil {
		ctxLogger.Error("failed to validate request body", err)

		Error(w, http.StatusBadRequest, "Failed validating request body.")

		return
	}

	ctxLogger = ctxLogger.With(
		slog.String("scooter_id", rentPost.ScooterUUID.String()),
		slog.String("city", rentPost.City),
	)

	rentalScooter := modelrental.NewRentInfo(rentPost.ScooterUUID.String(), rentPost.City)

	ctxLogger.Info("Renting scooter.")

	if err = s.rentalService.Rent(ctx, rentalScooter); err != nil {
		ctxLogger.Error("failed to rent a scooter", err)

		Error(w, http.StatusInternalServerError, "Failed renting scooter.")

		return
	}

	ctxLogger.Info("Successfully rented scooter.")

	trackerInfo := trackermodel.NewScooter(
		rentPost.ScooterUUID.String(),
		rentPost.City,
		rentPost.Longitude,
		rentPost.Latitude,
	)

	if err = s.trackerService.Track(clientUUID, trackerInfo); err != nil {
		ctxLogger.Warn("Failed to enable tracking for rented scooter.")
	} else {
		ctxLogger.Info("Tracking rented scooter.")
	}

	JSON(w, http.StatusNoContent, nil)
}

// freeScooter enables user to free the scooter that is used by the user.
//
//	@Summary	Free the given scooter.
//	@Tags		scooters
//
//	@Param		Client-Id	header	string			true	"ClientID"	minlength(36)	maxlength(36)	default(00000000-0000-0000-0000-000000000000)
//	@Param		Payload		body	model.FreePost	true	"Scooter to free information"
//
//	@Success	204
//	@Failure	400	{object}	model.ApiError
//	@Failure	403	{object}	model.ApiError
//	@Failure	500	{object}	model.ApiError
//	@Router		/free [post]
func (s *Server) freeScooter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clientUUID, err := clientUUIDFromHeader(r)
	if err != nil {
		s.logger.Error("failed to get clientID from header")

		Error(w, http.StatusBadRequest, "Failed getting clientUUID from header.")

		return
	}

	ctxLogger := s.logger.With(
		slog.String("client_id", clientUUID.String()),
	)

	var freePost model.FreePost

	if err = json.NewDecoder(r.Body).Decode(&freePost); err != nil {
		ctxLogger.Error("failed to decode scooterID from request body", err)

		Error(w, http.StatusBadRequest, "Failed decoding request body to scooterID.")

		return
	}

	if err = s.validator.Struct(freePost); err != nil {
		ctxLogger.Error("failed to validate request body", err)

		Error(w, http.StatusBadRequest, "Failed validating request body.")

		return
	}

	ctxLogger = s.logger.With(
		slog.String("scooter_id", freePost.ScooterUUID.String()),
	)

	ctxLogger.Info("Freeing the scooter.")

	if err = s.rentalService.Free(ctx, freePost.ScooterUUID); err != nil {
		ctxLogger.Error("failed to free the scooter", err)

		Error(w, http.StatusInternalServerError, "Failed freeing scooter.")

		return
	}

	ctxLogger.Info("Successfully freed the scooter.")

	if err = s.trackerService.StopTracking(clientUUID, freePost.ScooterUUID); err != nil {
		ctxLogger.Warn("Failed to stop tracking the scooter.", err)
	} else {
		ctxLogger.Info("Stopped tracking the scooter.")
	}

	JSON(w, http.StatusNoContent, nil)
}

func clientUUIDFromHeader(r *http.Request) (uuid.UUID, error) {
	clientUUIDAsString := r.Header.Get("Client-Id")
	if len(clientUUIDAsString) == 0 {
		return uuid.Nil, errExpectedHeaderParamNotFound
	}

	clientUUID, err := uuid.Parse(clientUUIDAsString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing clientID: %w", err)
	}

	return clientUUID, nil
}

// JSON writes a JSON response.
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set(headerContentType, contentTypeJSON)
	body, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"Error":"%v"}`, err.Error())))

		return
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
}

// Error writes an error response.
func Error(w http.ResponseWriter, statusCode int, message string) {
	apiError := model.ApiError{
		Message: message,
	}

	JSON(w, statusCode, &apiError)
}
