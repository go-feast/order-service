package order

import "github.com/google/uuid"

type Operation func(*Order) error

type Repository interface {
	Create(order *Order) error
	Get(id uuid.UUID) (*Order, error)
	Operate(id uuid.UUID, op Operation) error
}
