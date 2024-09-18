package model

type Scooter struct {
	Name                string
	Longitude, Latitude float64
	City                string
	Availability        bool
}

func NewScooter(name, city string, long, lat float64, availability bool) *Scooter {
	return &Scooter{
		Name:         name,
		City:         city,
		Longitude:    long,
		Latitude:     lat,
		Availability: availability,
	}
}
