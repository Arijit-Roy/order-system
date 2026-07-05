package repository

import (
	"context"
	"order-system/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOutboxRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOutboxRepository(pool *pgxpool.Pool) *PostgresOutboxRepository {
	return &PostgresOutboxRepository{
		pool: pool,
	}
}

func (r *PostgresOutboxRepository) ListPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error){
	rows, err := r.pool.Query(ctx, `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, status
		FROM outbox_events
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	events := make([]domain.OutboxEvent, 0, limit)

	for rows.Next() {
		var e domain.OutboxEvent
		if err := rows.Scan(
			&e.ID,
			&e.AggregateType,
			&e.AggregateID,
			&e.EventType,
			&e.Payload,
			&e.Status,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *PostgresOutboxRepository) MarkPublished(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE outbox_events
		SET status = 'PUBLISHED',
		    published_at = NOW()
		WHERE id = $1
	`, id)

	return err
}