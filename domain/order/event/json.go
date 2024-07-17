package event

import (
	"github.com/google/uuid"
	"service/domain/shared/destination"
	"service/event"
)

// JSONEventOrderCreated provides JSON representation of Order.
type JSONEventOrderCreated struct {
	event.Event  `json:"-"`
	OrderID      string                      `json:"order_id"`
	CustomerID   string                      `json:"customer_id"`
	RestaurantID string                      `json:"restaurant_id"`
	Meals        []string                    `json:"meals"`
	Destination  destination.JSONDestination `json:"destination"`
}

type JSONOrderFinished struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
}

type JSONOrderCooking struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
}

type JSONEventOrderPaid struct {
	event.Event   `json:"-"`
	OrderID       uuid.UUID `json:"order_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
}

type JSONWaitingForCourier struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
}

type JSONCourierTook struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
	CourierID   uuid.UUID `json:"courier_id"`
}

type JSONDelivering struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
}

type JSONDelivered struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
}

type JSONCanceled struct { //nolint:govet
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
	Reason      string    `json:"reason"`
}

type JSONClosed struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
}
