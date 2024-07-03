package order

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"service/domain/shared/destination"
	"time"
)

func InitializeOrderScheme(db *gorm.DB) {
	err := db.AutoMigrate(&DatabaseOrderDTO{})
	if err != nil {
		panic(errors.Wrap(err, "failed to migrate database"))
	}
}

type DatabaseOrderDTO struct { //nolint:govet
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey"`
	RestaurantID  uuid.UUID   `gorm:"type:uuid"`
	CustomerID    uuid.UUID   `gorm:"type:uuid"`
	CourierID     uuid.UUID   `gorm:"type:uuid"`
	Meals         []uuid.UUID `gorm:"type:uuid[]"`
	State         State       `gorm:"type:text"`
	TransactionID uuid.UUID   `gorm:"type:uuid"`
	Latitude      float64     `gorm:"type:numeric"`
	Longitude     float64     `gorm:"type:numeric"`
	CreatedAt     time.Time
}

func (d *DatabaseOrderDTO) ToOrder() *Order {
	dst, _ := destination.NewDestination(d.Latitude, d.Longitude)
	return &Order{
		id:            d.ID,
		restaurantID:  d.RestaurantID,
		customerID:    d.CustomerID,
		courierID:     d.CourierID,
		meals:         d.Meals,
		state:         d.State,
		transactionID: d.TransactionID,
		destination:   dst,
		createdAt:     d.CreatedAt,
	}
}

func (o *Order) ToDatabaseDTO() *DatabaseOrderDTO {
	return &DatabaseOrderDTO{
		ID:            o.id,
		RestaurantID:  o.restaurantID,
		CustomerID:    o.customerID,
		CourierID:     o.courierID,
		Meals:         o.meals,
		State:         o.state,
		TransactionID: o.transactionID,
		Latitude:      o.destination.Latitude(),
		Longitude:     o.destination.Longitude(),
		CreatedAt:     o.createdAt,
	}
}
