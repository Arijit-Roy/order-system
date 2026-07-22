package inventory

import (
	"context"
	"errors"
	"order-system/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) ReserveInventory(ctx context.Context, order domain.Order, event domain.OutboxEvent) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO orders (
			id,
			customer_id,
			product_id,
			quantity,
			amount,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (id)
		DO UPDATE SET
			customer_id = EXCLUDED.customer_id,
			product_id  = EXCLUDED.product_id,
			quantity    = EXCLUDED.quantity,
			amount      = EXCLUDED.amount,
			updated_at  = NOW()
		`,
		order.ID,
		order.CustomerID,
		order.ProductID,
		order.Quantity,
		order.Amount,
	)

	if err != nil {
		return err
	}

	cmd, err := tx.Exec(
		ctx,
		`UPDATE stock
		SET available_quantity = available_quantity - $1
		WHERE product_id = $2
		  AND available_quantity >= $1
		`,
		order.Quantity,
		order.ProductID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("Insufficient stock")
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO outbox_events
		 (id, aggregate_type, aggregate_id, event_type, payload, status)
		 VALUES ($1, $2, $3, $4, $5::jsonb, $6)`,
		event.ID,
		event.AggregateType,
		event.AggregateID,
		event.EventType,
		string(event.Payload),
		"PENDING",
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepo) UpsertOrder(ctx context.Context, order domain.Order) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO orders (
			id,
			customer_id,
			product_id,
			quantity,
			amount,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (id)
		DO UPDATE SET
			customer_id = EXCLUDED.customer_id,
			product_id  = EXCLUDED.product_id,
			quantity    = EXCLUDED.quantity,
			amount      = EXCLUDED.amount,
			currency    = EXCLUDED.currency,
			updated_at  = NOW()
		`,
		order.ID,
		order.CustomerID,
		order.ProductID,
		order.Quantity,
		order.Amount,
	)

	return err
}
