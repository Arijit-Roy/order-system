package repository

import (
	"context"
	"order-system/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type DBTX interface {
// 	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
// 	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
// }

type PostgresOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresOrderRepository{
	return &PostgresOrderRepository{
		pool: pool,
	}
}

func (r *PostgresOrderRepository) CreateOrder(
	ctx context.Context, 
	order domain.Order, 
	event domain.OutboxEvent,
) error{
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO orders (id, customer_id, product_id, quantity, amount)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (id) DO UPDATE
		 SET customer_id = EXCLUDED.customer_id,
		     product_id = EXCLUDED.product_id,
		     quantity = EXCLUDED.quantity,
		     amount = EXCLUDED.amount`,
		order.ID,
		order.CustomerID,
		order.ProductID,
		order.Quantity,
		order.Amount,
	)
	if err != nil {
		return err
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

func (r *PostgresOrderRepository) ListOrders(ctx context.Context) ([]domain.Order, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, customer_id, product_id, quantity, amount
		FROM orders
		ORDER BY created_at ASC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]domain.Order, 0)
	for rows.Next() {
		var o domain.Order
		err := rows.Scan(&o.ID, &o.CustomerID, &o.ProductID, &o.Quantity, &o.Amount)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

