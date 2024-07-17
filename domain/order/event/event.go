package event

import (
	"service/domain/shared/destination"
)

// EventType provides methods for converting Order for different marshaling strategies.
type EventType struct {
	OrderID       string
	CustomerID    string
	RestaurantID  string
	TransactionID string
	Meals         []string
	Destination   destination.Destination
}

// JSONEventOrderCreated converts EventType to JSONEventOrderCreated.
func (t *EventType) JSONEventOrderCreated() JSONEventOrderCreated {
	return JSONEventOrderCreated{
		OrderID:      t.OrderID,
		CustomerID:   t.CustomerID,
		RestaurantID: t.RestaurantID,
		Meals:        t.Meals,
		Destination:  t.Destination.ToJSON(),
	}
}
