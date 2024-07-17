package saver

import "context"

type Saver[T comparable] interface {
	Save(ctx context.Context, entity T) error
}
