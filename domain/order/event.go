package order

import (
	"service/domain/shared/destination"
	"service/eserializer"
)

type JSONEventOrderCreated struct {
	eserializer.Event
	OrderID      string                      `json:"order_id"`
	CustomerID   string                      `json:"customer_id"`
	RestaurantID string                      `json:"restaurant_id"`
	Meals        []string                    `json:"meals"`
	Destination  destination.JSONDestination `json:"destination"`
}
