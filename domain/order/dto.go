package order

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"service/domain/shared/destination"
	"time"
)

func InitializeOrderScheme(db *gorm.DB) {
	err := db.AutoMigrate(&DatabaseOrderDTO{}, &RestaurantOrderDTO{}, &Meal{})
	if err != nil {
		panic(errors.Wrap(err, "failed to migrate database"))
	}
}

type DatabaseOrderDTO struct { //nolint:govet
	ID            uuid.UUID          `gorm:"type:uuid;primaryKey"`
	RestaurantID  uuid.UUID          `gorm:"type:uuid"`
	Restaurant    RestaurantOrderDTO `gorm:"foreignKey:RestaurantID;references:ID"`
	CustomerID    uuid.UUID          `gorm:"type:uuid"`
	CourierID     uuid.UUID          `gorm:"type:uuid"`
	State         State              `gorm:"type:text"`
	TransactionID uuid.UUID          `gorm:"type:uuid"`
	Latitude      float64            `gorm:"type:numeric"`
	Longitude     float64            `gorm:"type:numeric"`
	CreatedAt     time.Time
}

type RestaurantOrderDTO struct { //nolint:govet
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Meals []Meal    `gorm:"foreignKey:RestaurantID;references:ID"`
}

type Meal struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:uuid"`
}

func (d *DatabaseOrderDTO) ToOrder() *Order {
	dst, _ := destination.NewDestination(d.Latitude, d.Longitude)

	meals := make(uuid.UUIDs, len(d.Restaurant.Meals))

	for i, meal := range d.Restaurant.Meals {
		meals[i] = meal.ID
	}

	return &Order{
		id:            d.ID,
		restaurantID:  d.Restaurant.ID,
		customerID:    d.CustomerID,
		courierID:     d.CourierID,
		meals:         meals,
		state:         d.State,
		transactionID: d.TransactionID,
		destination:   dst,
		createdAt:     d.CreatedAt,
	}
}

func (o *Order) ToDatabaseDTO() *DatabaseOrderDTO {
	meals := make([]Meal, len(o.meals))

	for i, meal := range o.meals {
		meals[i].ID = meal
	}

	return &DatabaseOrderDTO{
		ID: o.id,
		Restaurant: RestaurantOrderDTO{
			ID:    o.restaurantID,
			Meals: meals,
		},
		CustomerID:    o.customerID,
		CourierID:     o.courierID,
		State:         o.state,
		TransactionID: o.transactionID,
		Latitude:      o.destination.Latitude(),
		Longitude:     o.destination.Longitude(),
		CreatedAt:     o.createdAt,
	}
}