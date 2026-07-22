package notification

import (
	"context"
	"log"
	"order-system/internal/domain"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SendOrderConfirmation(ctx context.Context, order domain.Order) error {

	// send order email
	log.Printf(
		"Sending confirmation email for order %s to customer %s",
		order.ID,
		order.CustomerID,
	)

	return nil
}
