package model

const (
	Rent = iota
	Free
)

type RentInfo struct {
	ScooterUUID string
	City        string
}

func NewRentInfo(scooterUUID, city string) *RentInfo {
	return &RentInfo{
		ScooterUUID: scooterUUID,
		City:        city,
	}
}
