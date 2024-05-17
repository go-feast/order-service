package order

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type (
	// ID should be used for Order`s id.
	ID string

	// MealID should only stay for restaurant meals` id.
	MealID ID

	// Destination represents position coordinates.
	Destination struct {
		Latitude  float64
		Longitude float64
	}
)

// Order represents Order service domain.
// User order should be created by client and passed through network and deserialized into Order.
// Order must contain all fields to pass for a specific service for each.
type Order struct { //nolint:govet
	// ID states for Order id.
	ID ID

	// RestaurantID states for Restaurant id.
	RestaurantID ID

	// UserID states for User id.
	UserID ID

	// Meals contain meals` id that user have selected in a specific restaurant.
	Meals []MealID

	// Destination contains geo position of where Order should be delivered
	Destination Destination

	// CreatedAt represents where Order has been created.
	CreatedAt time.Time

	// Paid represents if order is paid
	Paid bool

	// Finished represents if order is finished
	Finished bool

	// Represents a timestamp when order was cooked in restaurant.
	FinishedAt *time.Time
}

// NewID parses provided id string and returning
func NewID(id string) (ID, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return ID(id), nil
}

type MealIDError struct {
	Err   error
	Index int
}

func (m MealIDError) Error() string {
	return fmt.Errorf("invalid id(index-%d): %w", m.Index, m.Err).Error()
}

func NewMealIDError(err error, index int) error {
	return &MealIDError{Err: err, Index: index}
}

func Meals(ids []string) ([]MealID, []error) {
	var (
		errs    = make([]error, 0, len(ids))
		mealIDs = make([]MealID, len(ids))
	)

	for i, id := range ids {
		newID, err := NewID(id)
		switch err {
		case nil:
			mealIDs[i] = MealID(newID)
		default:
			errs = append(errs, NewMealIDError(err, i))
		}
	}

	if len(errs) != 0 {
		return nil, errs
	}

	return mealIDs, nil
}
