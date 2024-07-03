package order

import (
	"context"
	"github.com/google/uuid"
)

type Operation func(*Order) error

type Repository interface {
	Create(ctx context.Context, order *Order) error
	Get(ctx context.Context, id uuid.UUID) (*Order, error)
	Operate(ctx context.Context, id uuid.UUID, op Operation) error
}
