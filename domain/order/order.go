// Package order contains our business logic.
package order

import (
	"github.com/google/uuid"
	"time"
)

type (
	ID string

	// Destination represents position coordinates.
	Destination struct {
		Latitude  float64
		Longitude float64
	}
)

func (id ID) GetID() string {
	return string(id)
}

// Order represents Order service domain.
// User order should be created by client and passed through network and deserialized into Order.
// Order must contain all fields to pass for a specific service for each.
type Order struct { //nolint:govet
	// ID states for Order id.
	ID ID

	// RestaurantID states for Restaurant id.
	RestaurantID ID

	// CustomerID states for User id.
	CustomerID ID

	// CourierID states for courier id.
	CourierID ID

	// Meals contain meals` id that user have selected in a specific restaurant.
	Meals []ID

	// Destination contains geo position of where Order should be delivered.
	Destination Destination

	// CreatedAt represents where Order has been created.
	CreatedAt time.Time

	// CashPayment represents if order will be paid in cash or card.
	CashPayment bool

	// Paid represents if payment successful.
	Paid bool

	// TransactionID represents payment transaction id.
	TransactionID *ID

	// Finished represents if order is finished.
	Finished bool

	// Represents a timestamp when order was cooked in restaurant.
	// By default - zeroed.
	FinishedAt time.Time
}

func (o *Order) GetMealsID() []string {
	arr := make([]string, len(o.Meals))

	for i, meal := range o.Meals {
		arr[i] = meal.GetID()
	}

	return arr
}

// IsFinished return true and time if orders when Finished set to true and FinishedAt not zero.
// Otherwise, returning false and zeroed time.
func (o *Order) IsFinished() (bool, time.Time) {
	if o.Finished && !o.FinishedAt.IsZero() {
		return true, o.FinishedAt
	}

	return false, time.Time{}
}

// Create assign order an id and initialize a CreatedAt field.
func (o *Order) Create() {
	o.CreatedAt = time.Now()
	o.ID = ID(uuid.NewString())
}

// NewOrder creates an order. But to set [Order.ID] and [Order.CreatedAt] call Create method.
func NewOrder(
	restaurantID, userID, transactionID string,
	mealsIDs []string,
	cashPayment bool,
	latitude, longitude float64,
) (*Order, error) {
	rid, err := NewID(restaurantID)
	if err != nil {
		return nil, ErrInvalidRestaurantID
	}

	uid, err := NewID(userID)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	meals, err := MealsID(mealsIDs)
	if err != nil {
		return nil, err
	}

	var (
		tid *ID
		id  ID
	)

	if !cashPayment {
		id, err = NewID(transactionID)
		if err != nil {
			return nil, ErrInvalidTransactionID
		}

		tid = &id
	}

	return &Order{
		RestaurantID:  rid,
		CustomerID:    uid,
		Meals:         meals,
		CashPayment:   cashPayment,
		Destination:   Destination{Latitude: latitude, Longitude: longitude},
		TransactionID: tid,
	}, nil
}

// NewID parses provided id string and returning
func NewID(id string) (ID, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return ID(id), nil
}

// MealsID convert provided ids in slice of MealID.
// If one error occurred while converting - an error returned.
func MealsID(ids []string) ([]ID, error) {
	var (
		errs    = make([]error, 0, len(ids))
		mealIDs = make([]ID, len(ids))
	)

	for i, id := range ids {
		newID, err := NewID(id)
		switch err {
		case nil:
			mealIDs[i] = newID
		default:
			errs = append(errs, NewMealIDError(err, i))
		}
	}

	if len(errs) != 0 {
		return nil, &MealsIDError{Errs: errs}
	}

	return mealIDs, nil
}
