package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/client/model"
	"github.com/google/uuid"
)

const (
	base         = "http://app:8081"
	api          = "/api"
	version      = "/v1"
	scootersPath = "/scooters"
	rentPath     = "/rent"
	freePath     = "/free"

	numberOfScooterRentals = 5
	timeOfScooterRentals   = 10
)

var errUnexpectedResponseStatus = errors.New("received unexpected response status")

type Customer struct {
	ClientUUID uuid.UUID
	Longitude  float64
	Latitude   float64
	Height     float64
	Width      float64
	City       string
}

type CustomerClient interface {
	SimulateUser(customer *Customer)
}

type customerClient struct {
	logger *slog.Logger
	client *http.Client
}

func NewCustomerClient(logger *slog.Logger) *customerClient {
	return &customerClient{
		logger: logger,
		client: http.DefaultClient,
	}
}

func (c *customerClient) SimulateUser(customer *Customer) {
	for i := 1; i <= numberOfScooterRentals; i++ {
		scooters, err := c.getScooters(customer)
		if err != nil {
			c.logger.Error("failed getting scooters information", err)

			return
		}

		availableScooters := filterAvailableScooters(scooters)
		if len(availableScooters) == 0 {
			c.logger.Info("No available scooters.")
			time.Sleep(time.Second)

			continue
		}

		j := rand.Intn(len(availableScooters))

		err = c.rentScooter(customer, &availableScooters[j], customer.City)
		if err != nil {
			c.logger.Error("failed renting scooter", err)
			i--

			time.Sleep(time.Second)

			continue
		}

		time.Sleep(timeOfScooterRentals * time.Second)

		err = c.freeScooter(customer, availableScooters[j].UUID)
		if err != nil {
			c.logger.Error("failed freeing scooter", err)

			return
		}
	}
}

func (c *customerClient) getScooters(client *Customer) ([]model.ScooterGet, error) {
	requestScooters, err := buildRequest(client, scootersPath, http.MethodGet, &bytes.Buffer{})
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	validURLQuery := &url.Values{}
	validURLQuery.Add("longitude", strconv.FormatFloat(client.Longitude, 'f', -1, 64))
	validURLQuery.Add("latitude", strconv.FormatFloat(client.Latitude, 'f', -1, 64))
	validURLQuery.Add("height", strconv.FormatFloat(client.Height, 'f', -1, 64))
	validURLQuery.Add("width", strconv.FormatFloat(client.Width, 'f', -1, 64))
	validURLQuery.Add("city", client.City)

	requestScooters.URL.RawQuery = validURLQuery.Encode()

	response, err := c.client.Do(requestScooters)
	if err != nil {
		return nil, fmt.Errorf("requesting scooters: %w", err)
	}

	defer func() {
		if err = response.Body.Close(); err != nil {
			c.logger.Error("failed closing body", err)

			return
		}
	}()

	switch response.StatusCode {
	case http.StatusOK:
		var scooters []model.ScooterGet
		if err = json.NewDecoder(response.Body).Decode(&scooters); err != nil {
			return nil, fmt.Errorf("decoding response body: %w", err)
		}

		return scooters, nil
	default:
		return nil, fmt.Errorf(
			"receiving response with status %s: %w",
			response.Status,
			errUnexpectedResponseStatus,
		)
	}
}

func (c *customerClient) rentScooter(client *Customer, scooter *model.ScooterGet, city string) error {
	scooterPost := model.ScooterPost{
		ScooterUUID:  scooter.UUID,
		Longitude:    scooter.Longitude,
		Latitude:     scooter.Latitude,
		Availability: scooter.Availability,
		City:         city,
	}

	scooterJSON, err := json.Marshal(scooterPost)
	if err != nil {
		return fmt.Errorf("marshaling scooter to JSON: %w", err)
	}

	requestRental, err := buildRequest(client, rentPath, http.MethodPost, bytes.NewBuffer(scooterJSON))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	response, err := c.client.Do(requestRental)
	if err != nil {
		return fmt.Errorf("requesting rental of a scooter: %w", err)
	}

	defer func() {
		if err = response.Body.Close(); err != nil {
			c.logger.Error("failed closing body", err)

			return
		}
	}()

	switch response.StatusCode {
	case http.StatusNoContent:
		c.logger.Info("Rented scooter successfully.", slog.String("scooter_id", scooter.UUID.String()))

		return nil
	default:
		return fmt.Errorf(
			"receiving response with status %s: %w",
			response.Status,
			errUnexpectedResponseStatus,
		)
	}
}

func (c *customerClient) freeScooter(client *Customer, scooterUUID uuid.UUID) error {
	scooter := model.FreePost{
		ScooterUUID: scooterUUID,
	}

	scooterJSON, err := json.Marshal(scooter)
	if err != nil {
		return fmt.Errorf("marshaling scooterUUID to JSON: %w", err)
	}

	requestFreeingScooter, err := buildRequest(client, freePath, http.MethodPost, bytes.NewBuffer(scooterJSON))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	response, err := c.client.Do(requestFreeingScooter)
	if err != nil {
		return fmt.Errorf("requesting freeing of the scooter: %w", err)
	}

	defer func() {
		if err = response.Body.Close(); err != nil {
			c.logger.Error("failed closing body", err)

			return
		}
	}()

	switch response.StatusCode {
	case http.StatusNoContent:
		c.logger.Info("Freed scooter successfully.", slog.String("scooter_id", scooterUUID.String()))

		return nil
	default:
		return fmt.Errorf(
			"receiving response with status %s: %w",
			response.Status,
			errUnexpectedResponseStatus,
		)
	}
}

func buildRequest(c *Customer, path string, method string, body *bytes.Buffer) (*http.Request, error) {

	request, err := http.NewRequestWithContext(
		context.Background(),
		method,
		base+api+version+path,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	request.Header.Set("Client-Id", c.ClientUUID.String())

	return request, nil
}

func filterAvailableScooters(scooters []model.ScooterGet) []model.ScooterGet {
	var filteredScooters []model.ScooterGet

	for i := range scooters {
		if scooters[i].Availability {
			filteredScooters = append(filteredScooters, scooters[i])
		}
	}

	return filteredScooters
}
