package service

import (
	"context"
	"encoding/json"
	"fmt"
	"order-system/internal/domain"
	"order-system/internal/repository"

	"github.com/google/uuid"
)

type OrderService struct {
	repo repository.OrderRepository
}

type orderCreatedPayload struct {
	ID         string  `json:"id"`
	CustomerID string  `json:"customer_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Amount     float64 `json:"amount"`
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService {
		repo: repo,
	}
}
func (s *OrderService) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	if(order.Amount<0){
		return domain.Order{}, fmt.Errorf("invalid amount");
	}
	payload, err := json.Marshal(orderCreatedPayload{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		ProductID:  order.ProductID,
		Quantity:   order.Quantity,
		Amount:     order.Amount,
	})

	if err != nil {
		return domain.Order{}, err
	}

	event := domain.OutboxEvent{
		ID:           uuid.NewString(),
		AggregateType: "order",
		AggregateID:   order.ID,
		EventType:     "order.created",
		Payload:       payload,
	}

	if err := s.repo.CreateOrder(ctx, order, event); err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func (s *OrderService) ListOrders(ctx context.Context) ([]domain.Order, error) {
	return s.repo.ListOrders(ctx)
}


func (s *OrderService) TotalRevenue(ctx context.Context) (float64, error) {
	orders, err := s.repo.ListOrders(ctx)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, o := range orders {
		total += o.Amount
	}

	return total, nil
}