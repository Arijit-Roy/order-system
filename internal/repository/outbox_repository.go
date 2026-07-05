package repository

import (
	"order-system/internal/domain"
)

type OutboxRepository interface {
	ListPending(ctx, limit int) ([]domain.OutboxEvent, error)
	MarkPublished(ctx, id string) error
}