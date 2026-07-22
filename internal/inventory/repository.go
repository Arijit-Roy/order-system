package inventory

import (
	"context"
	"order-system/internal/domain"
)

type Repository interface {
	UpsertOrder(ctx context.Context, order domain.Order) error
	ReserveInventory(ctx context.Context, order domain.Order, event domain.OutboxEvent) error
}
