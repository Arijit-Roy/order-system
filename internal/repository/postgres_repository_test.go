package repository

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"order-system/internal/domain"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func TestPostresRepo_SaveAndFindAll(t *testing.T){
// 	ctx,cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	pool, err := pgxpool.New(ctx, "postgres://orderuser:orderpass@localhost:5432/ordersdb")

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer pool.Close()
// 	tx, err := pool.Begin(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer tx.Rollback(ctx)

// 	repo := NewPostgresRepo(tx)

// 	order := domain.Order{
// 		ID:         uuid.NewString(),
// 		CustomerID: "CUST-001",
// 		ProductID:  "PROD-001",
// 		Quantity:   2,
// 		Amount:     100,
// 	}

// 	if err := repo.CreateOrder(ctx, order); err != nil {
// 		t.Fatal(err)
// 	}

// 	orders, err := repo.ListOrders(ctx)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(orders) == 0 {
// 		t.Fatal("expected at least one order")
// 	}

// 	found := false
// 	for _, o := range orders {
// 		if o.ID == order.ID {
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		t.Fatalf("saved order %s not found in query results", order.ID)
// 	}

// }

// func TestPostgresRepository_Save_WithCanceledContext(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
// 	defer cancel()

// 	pool, err := pgxpool.New(ctx, "postgres://orderuser:orderpass@localhost:5432/ordersdb")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer pool.Close()

// 	time.Sleep(1 * time.Millisecond)

// 	repo := NewPostgresRepo(pool)

// 	order := domain.Order{
// 		ID:         uuid.NewString(),
// 		CustomerID: "CUST-001",
// 		ProductID:  "PROD-001",
// 		Quantity:   2,
// 		Amount:     100,
// 	}

// 	err = repo.CreateOrder(ctx, order)

// 	if err == nil {
// 		t.Fatal("expected timeout/cancellation error, got nil")
// 	}

// 	if !errors.Is(err, context.DeadlineExceeded) {
// 		t.Fatalf("expected context deadline exceeded, got %v", err)
// 	}
// }

// func TestPostgresRepository_SaveDuplicateError(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	pool, err := pgxpool.New(ctx, "postgres://orderuser:orderpass@localhost:5432/ordersdb")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer pool.Close()

// 	repo := NewPostgresRepo(pool)

// 	id := uuid.NewString()

// 	order := domain.Order{
// 		ID:id,
// 		CustomerID: "CUST-001",
// 		ProductID:  "PROD-001",
// 		Quantity:   2,
// 		Amount:     100,
// 	}

// 	err = repo.CreateOrder(ctx, order)

// 	if err != nil {
// 		fmt.Printf("order was not saved %s", err)
// 		return;
// 	}

// 	err = repo.CreateOrder(ctx, order)
// 	if err == nil {
// 		t.Fatal("expected duplicate error, got nil")
// 	}

// 	if !errors.Is(err, domain.ErrDuplicateOrder) {
// 		t.Fatalf("expected ErrDuplicateOrder, got %v", err)
// 	}
// }

// func TestPostgresRepository_Save_WithExplicitlyCanceledContext(t *testing.T) {
// 	ctx := context.Background()

// 	pool, err := pgxpool.New(ctx, "postgres://orderuser:orderpass@localhost:5432/ordersdb")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer pool.Close()

// 	repo := NewPostgresRepo(pool)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	cancel()

// 	order := domain.Order{
// 		ID:         uuid.NewString(),
// 		CustomerID: "CUST-001",
// 		ProductID:  "PROD-001",
// 		Quantity:   2,
// 		Amount:     100,
// 	}

// 	err = repo.CreateOrder(ctx, order)
// 	if err == nil {
// 		t.Fatal("expected cancellation error, got nil")
// 	}
// 	if !errors.Is(err, context.Canceled) {
// 		t.Fatalf("expected context canceled, got %v", err)
// 	}
// }

// func TestPostgresRepository_Save_UpdatesExistingOrder(t *testing.T) {
// 	ctx := context.Background()

// 	pool, err := pgxpool.New(ctx, "postgres://orderuser:orderpass@localhost:5432/ordersdb")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer pool.Close()

// 	tx, err := pool.Begin(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer tx.Rollback(ctx)

// 	repo := NewPostgresRepo(tx)

// 	order := domain.Order{
// 		ID:         uuid.NewString(),
// 		CustomerID: "CUST-001",
// 		ProductID:  "PROD-001",
// 		Quantity:   2,
// 		Amount:     100,
// 	}

// 	if err := repo.CreateOrder(ctx, order); err != nil {
// 		t.Fatal(err)
// 	}

// 	order.Amount = 200

// 	if err := repo.CreateOrder(ctx, order); err != nil {
// 		t.Fatal(err)
// 	}

// 	orders, err := repo.ListOrders(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var found *domain.Order

// 	for _, o := range orders {
// 		if o.ID == order.ID {
// 			found = &o
// 		}
// 	}

// 	if found == nil {
// 		t.Fatalf("expected to find order %s", order.ID)
// 	}

// 	if found.Amount != 200 {
// 		t.Fatalf("expected amount 200, got %.2f", found.Amount)
// 	}
// }