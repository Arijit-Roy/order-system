package repository

import (
	"context"
	"fmt"
	"order-system/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order domain.Order, event domain.OutboxEvent) error
	ListOrders(ctx context.Context) ([]domain.Order, error)
}

// type InMemoryOrderRepository struct {
// 	orders []Order
// }

// func NewInMemoryOrderRepository() *InMemoryOrderRepository {
// 	return &InMemoryOrderRepository{
// 		orders: []Order{},
// 	}
// }

// func (r *InMemoryOrderRepository) Save(order Order) {
// 	r.orders = append(r.orders, order)
// }

// func (r *InMemoryOrderRepository) FindAll() []Order {
// 	return r.orders
// }





type FakePostgresOrderRepository struct {
	orders []domain.Order
}

func NewFakePostgresOrderRepository() *FakePostgresOrderRepository {
	return &FakePostgresOrderRepository{
		orders: []domain.Order{},
	}
}

func (r *FakePostgresOrderRepository) Save(order domain.Order) {
	fmt.Println("saving order to postgres")
	r.orders = append(r.orders, order)
}

func (r *FakePostgresOrderRepository) FindAll() []domain.Order {
	fmt.Println("listing order in postgres")
	return r.orders
}



// test repo
type MockOrderRepository struct {
	savedOrders []domain.Order
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		savedOrders: []domain.Order{},
	}
}

func (r *MockOrderRepository) Save(order domain.Order) error {
	r.savedOrders = append(r.savedOrders, order)
	return nil
}

func (r *MockOrderRepository) FindAll() []domain.Order {
	return r.savedOrders
}