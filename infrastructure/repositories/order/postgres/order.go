package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"service/domain/order"
)

const (
	selectState = "select orders.state from orders where id = :id"
)

type PostgreSQLRepository struct { //nolint:revive
	execer   sqlx.ExecerContext
	queryer  sqlx.QueryerContext
	preparer sqlx.PreparerContext
	db       *sqlx.DB
}

func NewPostgreSQLRepository(db *sqlx.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		execer:   db,
		queryer:  db,
		preparer: db,
		db:       db,
	}
}

var opts = sql.TxOptions{Isolation: sql.LevelDefault} //nolint:unused

func (p PostgreSQLRepository) executeTx(ctx context.Context, operation func(tx *sqlx.Tx) error) (err error) {
	txx, err := p.db.BeginTxx(ctx /*opts*/, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create postgres transaction")
	}

	defer func() {
		if v := recover(); v != nil {
			// err should be always nil here
			err = multierror.Prefix(fmt.Errorf("panic occurred: %s", v), "failed to process transaction: ")
		}

		if err != nil {
			e := txx.Rollback()
			if e != nil {
				err = multierror.Append(err, e)
				err = errors.Wrap(err, "failed to rollback transaction")

				return
			}
		}
	}()

	err = operation(txx)
	if err != nil {
		return errors.Wrap(err, "failed to perform transaction script")
	}

	err = txx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (p PostgreSQLRepository) Create(ctx context.Context, _ *order.Order) error {
	return p.executeTx(ctx, func(_ *sqlx.Tx) error {
		panic("implement me")
	})
}

func (p PostgreSQLRepository) Get(_ context.Context, _ uuid.UUID) (*order.Order, error) {
	panic("implement me")
}

func (p PostgreSQLRepository) Operate(ctx context.Context, id uuid.UUID, op order.Operation) error {
	return p.executeTx(ctx, func(tx *sqlx.Tx) error {
		o, err := p.WithTx(tx).Get(ctx, id)
		if err != nil {
			return errors.Wrap(err, "failed to get order")
		}

		var v any
		// select order state for preventing data race
		err = tx.SelectContext(ctx, v, selectState, id)
		if err != nil {
			return errors.Wrap(err, "failed to select order`s state in transaction")
		}

		return op(o)
	})
}

func (p PostgreSQLRepository) WithTx(tx *sqlx.Tx) PostgreSQLRepository {
	return PostgreSQLRepository{
		execer:   tx,
		queryer:  p.queryer,
		preparer: tx,
		db:       p.db,
	}
}
