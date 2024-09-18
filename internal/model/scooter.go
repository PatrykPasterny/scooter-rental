package model

import (
	"github.com/google/uuid"
)

type ScooterGet struct {
	ScooterUUID  uuid.UUID `json:"UUID"`
	Longitude    float64   `json:"longitude"`
	Latitude     float64   `json:"latitude"`
	Availability bool      `json:"availability"`
}

type RentPost struct {
	ScooterUUID uuid.UUID `json:"UUID" validate:"required"`
	Longitude   float64   `json:"longitude" validate:"required"`
	Latitude    float64   `json:"latitude" validate:"required"`
	City        string    `json:"city" validate:"required"`
}

type FreePost struct {
	ScooterUUID uuid.UUID `json:"UUID" validate:"required"`
}

type ApiError struct {
	Message string `json:"Message"`
}
