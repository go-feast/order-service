package destination

type JSONDestination struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

func (d Destination) ToJSON() JSONDestination {
	return JSONDestination{
		Latitude:  d.latitude,
		Longitude: d.longitude,
	}
}
