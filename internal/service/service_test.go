package service

// import "testing"

// func TestCreateOrder(t *testing.T){
// 	repo := NewMockOrderRepository()
// 	service := NewOrderService(repo)

// 	_, err := service.CreateOrder(Order{
// 		ID:     "ORD-001",
// 		Amount: 100,
// 	})
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if len(repo.savedOrders) != 1 {
// 		t.Fatalf(
// 			"expected 1 order, got %d",
// 			len(repo.savedOrders),
// 		)
// 	}
// }

// func TestCreateOrder_Invalid(t *testing.T){
// 	repo := NewMockOrderRepository()
// 	service := NewOrderService(repo)

// 	_, err := service.CreateOrder(Order{
// 		ID:     "ORD-001",
// 		Amount: -100,
// 	})
// 	if err == nil {
// 		t.Fatal("expected error")
// 	}

// 	if len(repo.savedOrders) != 0 {
// 		t.Fatal("invalid order should not be saved")
// 	}
// }