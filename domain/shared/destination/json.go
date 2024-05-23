package destination

import "service/eserializer"

type JSONDestination struct {
	eserializer.Event
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (d Destination) ToJSON() JSONDestination {
	return JSONDestination{
		Latitude:  d.latitude,
		Longitude: d.longitude,
	}
}
