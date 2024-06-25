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
	selectOrder = `SELECT
    				id,
					restaurant_id,
					customer_id,
					courier_id,
					meals,
					state,
					transaction_id,
					destination,
					created_at
				FROM orders
				WHERE 
				    id = $1;`
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
	updateOrder = `
		UPDATE orders
		SET
		    restaurant_id = :restaurant_id,
		    customer_id = :customer_id,
		    courier_id = :courier_id,
		    meals = :meals,
		    state = :state,
		    destination = ST_GeomFromText(:destination, 4326),
		    transaction_id = :transaction_id
		WHERE
		    id = :id;
	`
)

type PostgreSQLRepository struct { //nolint:revive
	execer   sqlx.ExecerContext
	queryer  sqlx.QueryerContext
	preparer sqlx.PreparerContext
	db       *sqlx.DB
}

func (p PostgreSQLRepository) Create(ctx context.Context, o *order.Order) error {
	return p.Insert(ctx, o)
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
	ErrFailedOperation   = errors.New("failed operation")
)

func (p PostgreSQLRepository) executeTx(
	ctx context.Context,
	operation func(
		ctx context.Context,
		stmt *sqlx.Tx,
	) error,
) (err error) {
	txx, err := p.db.BeginTxx(ctx, &opts /*opts*/)
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

	err = operation(opCtx, txx)
	if err != nil {
		return errors.Wrap(err, "failed to perform transaction script")
	}

	err = txx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (p PostgreSQLRepository) Insert(ctx context.Context, o *order.Order) error {
	dto := o.ToDto()

	_, err := p.execer.ExecContext(ctx, insertOrder,
		dto.ID,
		dto.RestaurantID,
		dto.CustomerID,
		dto.CourierID,
		dto.Meals,
		dto.State,
		dto.TransactionID,
		dto.Destination.Longitude, //for postgis geography longitude is firsy to fill
		dto.Destination.Latitude,
		dto.CreatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert script")
	}

	return nil
}

func (p PostgreSQLRepository) Get(ctx context.Context, id uuid.UUID) (o *order.Order, err error) {
	err = p.executeTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		stmt, err := tx.PreparexContext(ctx, selectOrder)
		if err != nil {
			return errors.Wrap(err, "failed to prepare statement")
		}

		dto := &order.DatabaseOrderDTO{}
		err = stmt.GetContext(ctx, dto, id)
		if err != nil {
			return errors.Wrap(err, "failed to perform select statement")
		}

		if dto == nil {
			return errors.Wrap(ErrFailedOperation, "dto is nil")
		}

		o = dto.ToOrder()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (p PostgreSQLRepository) Operate(ctx context.Context, id uuid.UUID, op order.Operation) error {
	return p.executeTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// prepare
		selectStmt, err := tx.PreparexContext(ctx, selectOrder)
		if err != nil {
			return errors.Wrap(err, "failed to prepare select statement")
		}
		defer selectStmt.Close()

		// select
		dto := &order.DatabaseOrderDTO{}
		err = selectStmt.GetContext(ctx, dto, id)
		if err != nil {
			return errors.Wrap(err, "failed to select order")
		}

		// op
		o := dto.ToOrder()

		err = op(o)
		if err != nil {
			return errors.Wrap(err, "failed to perform operation with order")
		}

		// prepare
		updateStmt, err := tx.PrepareNamedContext(ctx, updateOrder)
		if err != nil {
			return errors.Wrap(err, "failed to prepare insert statement")
		}
		defer updateStmt.Close()

		// update
		_, err = updateStmt.ExecContext(ctx, o.ToDto())
		if err != nil {
			return errors.Wrap(err, "failed to perform update statement")
		}

		return nil
	})
}
