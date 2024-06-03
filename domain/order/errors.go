package order

import (
	"errors"
)

var (
	ErrInvalidState  = errors.New("invalid state")
	ErrOrderClosed   = errors.New("order closed")
	ErrOrderCanceled = errors.New("order canceled")
)
