// Package order contains our business logic.
package order

import (
	"github.com/google/uuid"
	"service/domain/shared/destination"
	"time"
)

// Order represents Order service domain.
// Should be created by client and passed through network and deserialized into Order.
type Order struct { //nolint:govet
	// id states for Order [uuid].
	id uuid.UUID

	// restaurantID states for restaurant [uuid].
	restaurantID uuid.UUID

	// customerID states for user [uuid].
	customerID uuid.UUID

	// courierID states for courier [uuid].
	courierID uuid.UUID //nolint:unused

	// meals states for meals` [uuid] that user selected in a specific restaurant.
	meals uuid.UUIDs

	// state states for Order State.
	//
	// Every State can go into Canceled State. But the only way where Canceled can go into is Closed.
	// Canceled -> Closed
	//
	// State machine for an order:
	// Created -> Paid -> Cooking -> Finished -> WaitingForCourier -> CourierTook -> Delivering -> Delivered -> Closed.
	//
	state State

	// transactionID represents payment transaction [uuid].
	transactionID uuid.UUID

	// destination contains geo position of where Order should be delivered.
	destination destination.Destination

	// createdAt represents where Order has been created.
	createdAt time.Time
}

// IsCanceled shows if order has been canceled.
func (o *Order) IsCanceled() bool {
	return o.state == Canceled
}

// IsClosed shows if order has been closed.
func (o *Order) IsClosed() bool {
	return o.state == Closed
}

// ToEvent converts Order to EventType.
func (o *Order) ToEvent() *EventType {
	return &EventType{
		OrderID:      o.id.String(),
		CustomerID:   o.customerID.String(),
		RestaurantID: o.restaurantID.String(),
		Meals:        o.meals.Strings(),
		Destination:  o.destination,
	}
}

// NewOrder creates new Order.
func NewOrder(
	restaurantID, userID, transactionID string,
	mealsIDs []string,
	latitude, longitude float64,
) (*Order, error) {
	errs := make([]error, 0)

	rid, err := uuid.Parse(restaurantID)
	if err != nil {
		errs = append(errs, ErrInvalidRestaurantID)
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		errs = append(errs, ErrInvalidUserID)
	}

	meals, err := mealsID(mealsIDs)
	if err != nil {
		errs = append(errs, err)
	}

	tid, err := uuid.Parse(transactionID)
	if err != nil {
		errs = append(errs, ErrInvalidTransactionID)
	}

	deliverTo, err := destination.NewDestination(latitude, longitude)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return nil, NewMultipleErrors("failed order validation", errs)
	}

	return &Order{
		id:           uuid.New(),
		restaurantID: rid,
		customerID:   uid,
		// courierID:
		meals:         meals,
		state:         Created,
		transactionID: tid,
		destination:   deliverTo,
		createdAt:     time.Now(),
	}, nil
}

// mealsID convert provided ids in slice of MealID.
// If one error occurred while converting - an error returned.
func mealsID(ids []string) (uuid.UUIDs, error) {
	var (
		errs    = make([]error, 0, len(ids))
		mealIDs = make(uuid.UUIDs, len(ids))
	)

	for i, id := range ids {
		newID, err := uuid.Parse(id)
		switch err {
		case nil:
			mealIDs[i] = newID
		default:
			errs = append(errs, NewMealIDError(err, i))
		}
	}

	if len(errs) != 0 {
		return nil, NewMultipleErrors("meals id error", errs)
	}

	return mealIDs, nil
}
