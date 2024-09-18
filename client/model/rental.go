package model

import "github.com/google/uuid"

type ScooterGet struct {
	UUID         uuid.UUID `json:"UUID"`
	Longitude    float64   `json:"longitude"`
	Latitude     float64   `json:"latitude"`
	Availability bool      `json:"availability"`
}

type ScooterPost struct {
	ScooterUUID  uuid.UUID `json:"UUID"`
	Longitude    float64   `json:"longitude"`
	Latitude     float64   `json:"latitude"`
	Availability bool      `json:"availability"`
	City         string    `json:"city"`
}

type RentPost struct {
	ScooterUUID uuid.UUID `json:"UUID"`
	Longitude   float64   `json:"longitude"`
	Latitude    float64   `json:"latitude"`
	City        string    `json:"city"`
}

type FreePost struct {
	ScooterUUID uuid.UUID `json:"UUID"`
}
