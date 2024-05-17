package order

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidRestaurantID  = errors.New("invalid restaurant id")
	ErrInvalidUserID        = errors.New("invalid user id")
	ErrInvalidTransactionID = errors.New("invalid transaction id")
)

type MealIDError struct {
	Err   error
	Index int
}

type MealsIDError struct {
	Errs []error `json:"errs"`
}

func (m *MealsIDError) Error() string {
	var sb strings.Builder

	sb.WriteString("meals ID errors: [")

	for i, err := range m.Errs {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(fmt.Sprintf("%v", err))
	}

	sb.WriteString("]")

	return sb.String()
}

func (m MealIDError) Error() string {
	return fmt.Errorf("invalid id(index=%d): %w", m.Index, m.Err).Error()
}

func NewMealIDError(err error, index int) error {
	return &MealIDError{Err: err, Index: index}
}
