package model

type ScooterQueryParams struct {
	Longitude    float64 `json:"longitude" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required"`
	Height       float64 `json:"height" validate:"required"`
	Width        float64 `json:"width" validate:"required"`
	City         string  `json:"city" validate:"required"`
	Availability *bool   `json:"availability"`
}
