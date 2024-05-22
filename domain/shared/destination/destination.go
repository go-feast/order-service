package destination

import "errors"

// Destination represents position coordinates.
type Destination struct {
	latitude  float64
	longitude float64
}

func (d Destination) Latitude() float64 {
	return d.latitude
}

func (d Destination) Longitude() float64 {
	return d.latitude
}

var (
	ErrInvalidLatitude  = errors.New("invalid latitude: must be between -90 and 90")
	ErrInvalidLongitude = errors.New("invalid longitude: must be between -180 and 180")
)

func NewDestination(lat, long float64) (Destination, error) {
	if !(lat >= -90 && lat <= 90) {
		return Destination{}, ErrInvalidLatitude
	}

	if !(long >= -180 && long <= 180) {
		return Destination{}, ErrInvalidLongitude
	}

	return Destination{
		latitude:  lat,
		longitude: long,
	}, nil
}
