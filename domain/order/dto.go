package order

import (
	"github.com/google/uuid"
	"service/domain/shared/destination"
	"time"
)

type DatabaseOrderDTO struct { //nolint:govet
	ID            uuid.UUID               `db:"id"`
	RestaurantID  uuid.UUID               `db:"restaurant_id"`
	CustomerID    uuid.UUID               `db:"customer_id"`
	CourierID     uuid.NullUUID           `db:"courier_id"`
	Meals         uuid.UUIDs              `db:"meals"`
	State         State                   `db:"state"`
	TransactionID uuid.NullUUID           `db:"transaction_id"`
	Destination   destination.Destination `db:"destination"`
	CreatedAt     time.Time               `db:"created_at"`
}

func (o *Order) ToDto() *DatabaseOrderDTO {
	return &DatabaseOrderDTO{
		ID:            o.id,
		RestaurantID:  o.restaurantID,
		CustomerID:    o.customerID,
		CourierID:     uuid.NullUUID{UUID: o.courierID, Valid: o.courierID != uuid.Nil},
		Meals:         o.meals,
		State:         o.state,
		TransactionID: uuid.NullUUID{UUID: o.transactionID, Valid: o.transactionID != uuid.Nil},
		Destination:   o.destination,
		CreatedAt:     o.createdAt,
	}
}

func (d *DatabaseOrderDTO) ToOrder() *Order {
	return &Order{
		id:            d.ID,
		restaurantID:  d.RestaurantID,
		customerID:    d.CustomerID,
		courierID:     d.CourierID.UUID,
		meals:         d.Meals,
		state:         d.State,
		transactionID: d.TransactionID.UUID,
		destination:   d.Destination,
		createdAt:     d.CreatedAt,
	}
}
