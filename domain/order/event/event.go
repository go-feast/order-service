package event

import (
	"service/domain/shared/destination"
)

// Type provides methods for converting Order for different marshaling strategies.
type Type struct {
	OrderID       string
	CustomerID    string
	RestaurantID  string
	TransactionID string
	Meals         []string
	Destination   destination.Destination
}

// JSONEventOrderCreated converts Type to JSONEventOrderCreated.
func (t *Type) JSONEventOrderCreated() JSONEventOrderCreated {
	return JSONEventOrderCreated{
		OrderID:      t.OrderID,
		CustomerID:   t.CustomerID,
		RestaurantID: t.RestaurantID,
		Meals:        t.Meals,
		Destination:  t.Destination.ToJSON(),
	}
}
