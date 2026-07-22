package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"order-system/internal/domain"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) HandleOrderCreated(ctx context.Context, order domain.Order) error {
	return s.repo.UpsertOrder(ctx, order)
}

func (s *Service) Reserve(ctx context.Context, order domain.Order) error {
	// if err := s.repo.UpsertOrder(ctx, order); err != nil {
	// 	return err
	// }
	payload, err := json.Marshal(order)
	if err != nil {
		return err
	}
	event := domain.OutboxEvent{
		ID:            uuid.NewString(),
		AggregateType: "inventory",
		EventType:     "inventory.reserved",
		Payload:       payload,
	}

	if err := s.repo.ReserveInventory(ctx, order, event); err != nil {
		fmt.Println("errrror is", err, order)
		return err
	}

	log.Printf("Reserved %d units for Product %s", order.Quantity, order.ProductID)
	return nil
}
