package model

type Scooter struct {
	Name                string
	Longitude, Latitude float64
	City                string
}

func NewScooter(name, city string, long, lat float64) *Scooter {
	return &Scooter{
		Name:      name,
		Longitude: long,
		Latitude:  lat,
		City:      city,
	}
}
