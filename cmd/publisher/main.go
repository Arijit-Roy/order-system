package main

import (
	"context"
	"fmt"
	"log"
	"order-system/internal/repository"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main(){
	ctx := context.Background();

	pool, err := pgxpool.New(ctx, "postgres://orderuser:orderpass@localhost:5432/ordersdb")

	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	outboxRepo := repository.NewPostgresOutboxRepository(pool)

	ticker := time.NewTicker(5*time.Second)
	defer ticker.Stop()

	log.Println("outbox publisher started")

	for {
		events, err := outboxRepo.ListPending(ctx, 100)
		if err != nil {
			log.Println("list pending failed:", err)
			<-ticker.C
			continue
		}

		for _, e := range events {
			fmt.Printf("publishing event %s type=%s aggregate=%s payload=%s\n",
				e.ID, e.EventType, e.AggregateID, string(e.Payload))

			if err := outboxRepo.MarkPublished(ctx, e.ID); err != nil {
				log.Println("mark published failed:", err)
			}
		}
		<-ticker.C
	}
}