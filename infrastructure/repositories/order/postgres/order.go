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
	"time"
)

const (
	selectState = "select orders.state from orders where id = :id"
	insertOrder = `INSERT INTO orders (
					id,
					restaurant_id,
					customer_id,
					courier_id,
					meals,
					state,
					transaction_id,
					destination,
					created_at
				) VALUES (
					$1,
					$2,
					$3,
					$4,
					$5,
					$6,
					$7,
					ST_SetSRID(ST_MakePoint($8, $9), 4326),
					$10
				);`
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

var (
	opts           = sql.TxOptions{Isolation: sql.LevelDefault} //nolint:unused
	defaultTimeout = 30 * time.Second
)

var (
	ErrChangesNotApplied = errors.New("changes not applied")
)

func (p PostgreSQLRepository) executeTx(ctx context.Context, operation func(ctx context.Context, p sqlx.PreparerContext) error) (err error) {
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

	opCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	operation(opCtx, txx)

	if err != nil {
		return errors.Wrap(err, "failed to perform transaction script")
	}

	err = txx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (p PostgreSQLRepository) Insert(ctx context.Context, order *order.Order) error {
	panic("implement me")
}

func (p PostgreSQLRepository) Create(ctx context.Context, o *order.Order) error {
	return p.executeTx(ctx, func(ctx context.Context, preparer sqlx.PreparerContext) error {
		script, err := preparer.PrepareContext(ctx, insertOrder)
		if err != nil {
			return errors.Wrap(err, "failed to prepare insert script")
		}

		dto := o.ToDto()

		result, err := script.ExecContext(ctx,
			dto.ID,
			dto.RestaurantID,
			dto.CustomerID,
			dto.CourierID,
			dto.Meals,
			dto.State,
			dto.TransactionID,
			dto.Destination.Longitude(), //for postgis geography longitude is firsy to fill
			dto.Destination.Latitude(),
			dto.CreatedAt,
		)
		if err != nil {
			return errors.Wrap(err, "failed to execute insert script")
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affected == 0 {
			return ErrChangesNotApplied
		}

		return nil
	})
}

func (p PostgreSQLRepository) Get(_ context.Context, _ uuid.UUID) (*order.Order, error) {
	panic("implement me")
}

func (p PostgreSQLRepository) Operate(ctx context.Context, id uuid.UUID, op order.Operation) error {
	return p.executeTx(ctx, func(ctx context.Context, preparer sqlx.PreparerContext) error {
		o, err := p.WithPreparer(preparer).Get(ctx, id)
		if err != nil {
			return errors.Wrap(err, "failed to get order")
		}

		if err = op(o); err != nil {
			return errors.Wrap(err, "failed to execute operation")
		}

		return p.WithPreparer(preparer).Insert(ctx, o)
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

func (p PostgreSQLRepository) WithPreparer(preparer sqlx.PreparerContext) PostgreSQLRepository {
	return PostgreSQLRepository{
		execer:   p.execer,
		queryer:  p.queryer,
		preparer: preparer,
		db:       p.db,
	}
}
