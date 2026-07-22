package publisher

import (
	"context"
	"order-system/internal/domain"
)

type Repository interface {
	ListPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	MarkPublished(ctx context.Context, id string) error
}
