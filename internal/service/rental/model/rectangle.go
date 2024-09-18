package model

type GeoRectangle struct {
	City                                           string
	CenterLongitude, CenterLatitude, Height, Width float64
}

func NewRectangle(city string, long, lat, height, width float64) *GeoRectangle {
	return &GeoRectangle{
		City:            city,
		CenterLongitude: long,
		CenterLatitude:  lat,
		Height:          height,
		Width:           width,
	}
}
