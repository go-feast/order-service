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
	ErrOrderClosed          = errors.New("order closed")
	ErrCannotSetState       = errors.New("cannot set state")
	ErrNoSuchState          = errors.New("no such state")
	ErrCannotSetPastState   = errors.New("cannot set past state")
)

type MealIDError struct {
	Err   error
	Index int
}

type MultipleErrors struct {
	Prefix string  `json:",omitempty"`
	Errs   []error `json:"errs"`
}

func NewMultipleErrors(prefix string, errs []error) *MultipleErrors {
	return &MultipleErrors{Prefix: prefix, Errs: errs}
}

func (m *MultipleErrors) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s: [", m.Prefix))

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
