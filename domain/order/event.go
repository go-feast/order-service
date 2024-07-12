package order

import (
	"github.com/google/uuid"
	"service/domain/shared/destination"
	"service/event"
)

// EventType provides methods for converting Order for different marshaling strategies.
type EventType struct {
	OrderID      string
	CustomerID   string
	RestaurantID string
	Meals        []string
	Destination  destination.Destination
}

// ToJSON converts EventType to JSONEventOrderCreated.
func (t *EventType) ToJSON() JSONEventOrderCreated {
	return JSONEventOrderCreated{
		OrderID:      t.OrderID,
		CustomerID:   t.CustomerID,
		RestaurantID: t.RestaurantID,
		Meals:        t.Meals,
		Destination:  t.Destination.ToJSON(),
	}
}

// JSONEventOrderCreated provides JSON representation of Order.
type JSONEventOrderCreated struct {
	event.Event
	OrderID      string                      `json:"order_id"`
	CustomerID   string                      `json:"customer_id"`
	RestaurantID string                      `json:"restaurant_id"`
	Meals        []string                    `json:"meals"`
	Destination  destination.JSONDestination `json:"destination"`
}

type JSONEventOrderPaid struct {
	event.Event
	OrderID       uuid.UUID `json:"order_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
}
