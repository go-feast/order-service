package order

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidState  = errors.New("invalid state")
	ErrOrderClosed   = errors.New("order closed")
	ErrOrderCanceled = errors.New("order setted")
)

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

		sb.WriteString(fmt.Sprintf("%s", err))
	}

	sb.WriteString("]")

	return sb.String()
}
