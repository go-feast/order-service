package order

import (
	"github.com/google/uuid"
	"service/domain/shared/destination"
	"time"
)

type DatabaseOrderDTO struct { //nolint:govet
	ID            uuid.UUID     `db:"id"`
	RestaurantID  uuid.UUID     `db:"restaurant_id"`
	CustomerID    uuid.UUID     `db:"customer_id"`
	CourierID     uuid.NullUUID `db:"courier_id"`
	Meals         []string      `db:"meals"`
	State         State         `db:"state"`
	TransactionID uuid.NullUUID `db:"transaction_id"`
	Destination   dst           `db:"destination"`
	CreatedAt     time.Time     `db:"created_at"`
}

type dst struct {
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}

func (o *Order) ToDto() *DatabaseOrderDTO {
	return &DatabaseOrderDTO{
		ID:            o.id,
		RestaurantID:  o.restaurantID,
		CustomerID:    o.customerID,
		CourierID:     uuid.NullUUID{UUID: o.courierID, Valid: true},
		Meals:         o.meals.Strings(),
		State:         o.state,
		TransactionID: uuid.NullUUID{UUID: o.transactionID, Valid: true},
		Destination:   dst{Latitude: o.destination.Latitude(), Longitude: o.destination.Longitude()},
		CreatedAt:     o.createdAt,
	}
}

func (d *DatabaseOrderDTO) ToOrder() *Order {
	point, _ := destination.NewDestination(d.Destination.Latitude, d.Destination.Longitude)
	meals, _ := mealsID(d.Meals)
	return &Order{
		id:            d.ID,
		restaurantID:  d.RestaurantID,
		customerID:    d.CustomerID,
		courierID:     d.CourierID.UUID,
		meals:         meals,
		state:         d.State,
		transactionID: d.TransactionID.UUID,
		destination:   point,
		createdAt:     d.CreatedAt,
	}
}
