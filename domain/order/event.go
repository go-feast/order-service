package order

import (
	"service/domain/shared/destination"
	"service/eserializer"
)

type EventType struct {
	OrderID      string
	CustomerID   string
	RestaurantID string
	Meals        []string
	Destination  destination.Destination
}

func (t *EventType) ToJSON() JSONEventOrderCreated {
	return JSONEventOrderCreated{
		OrderID:      t.OrderID,
		CustomerID:   t.CustomerID,
		RestaurantID: t.RestaurantID,
		Meals:        t.Meals,
		Destination:  t.Destination.ToJSON(),
	}
}

type JSONEventOrderCreated struct {
	eserializer.Event
	OrderID      string                      `json:"order_id"`
	CustomerID   string                      `json:"customer_id"`
	RestaurantID string                      `json:"restaurant_id"`
	Meals        []string                    `json:"meals"`
	Destination  destination.JSONDestination `json:"destination"`
}
