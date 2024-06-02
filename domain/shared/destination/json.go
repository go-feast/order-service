package destination

import "service/event"

type JSONDestination struct {
	event.Event
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (d Destination) ToJSON() JSONDestination {
	return JSONDestination{
		Latitude:  d.latitude,
		Longitude: d.longitude,
	}
}
