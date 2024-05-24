// Package order contains our business logic.
package order

import (
	"github.com/google/uuid"
	"service/domain/shared/destination"
	"time"
)

// Order represents Order service domain.
// User order should be created by client and passed through network and deserialized into Order.
// Order must contain all fields to pass for a specific service for each.
type Order struct { //nolint:govet
	// ID states for Order id.
	id uuid.UUID

	// rid states for Restaurant [uuid].
	restaurantID uuid.UUID

	// CustomerID states for User [uuid].
	customerID uuid.UUID

	// CourierID states for courier [uuid].
	courierID uuid.UUID //nolint:unused

	// meals states for meals` [uuid] that user have selected in a specific restaurant.
	meals uuid.UUIDs

	// state states for Order State. It could be one of
	//
	// Every State can go into Canceled State.
	// Canceled -> Closed
	//
	// Created -> Paid -> Cooking -> Finished -> WaitingForCourier -> CourierTook -> Delivering -> Delivered -> Closed.
	//
	state State

	// TransactionID represents payment transaction [uuid].
	transactionID uuid.UUID

	// Destination contains geo position of where Order should be delivered.
	destination destination.Destination

	// CreatedAt represents where Order has been created.
	createdAt time.Time
}

func (o *Order) IsCanceled() bool {
	return o.state == Canceled
}

func (o *Order) IsClosed() bool {
	return o.state == Closed
}

func (o *Order) ToEvent() *EventType {
	return &EventType{
		OrderID:      o.id.String(),
		CustomerID:   o.customerID.String(),
		RestaurantID: o.restaurantID.String(),
		Meals:        o.meals.Strings(),
		Destination:  o.destination,
	}
}

// NewOrder creates an order. But to set [Order.ID] and [Order.CreatedAt] call Create method.
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
