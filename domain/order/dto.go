package order

import (
	"database/sql"
	"github.com/google/uuid"
	"service/domain/shared/destination"
	"time"
)

type DatabaseOrderDTO struct { //nolint:govet
	ID            uuid.UUID
	RestaurantID  uuid.UUID
	CustomerID    uuid.UUID
	CourierID     uuid.NullUUID
	Meals         uuid.UUIDs
	State         State
	TransactionID uuid.NullUUID
	Destination   destination.Destination
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
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
