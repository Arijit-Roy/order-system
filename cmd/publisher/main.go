package main

import (
	"context"
	"log"
	"order-system/internal/config"
	"order-system/internal/publisher"
	"order-system/internal/redis"
	"order-system/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	pgDSN := cfg.PostgresDSN()

	orderDB, err := pgxpool.New(ctx, pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	orderRepo := repository.NewPostgresOutboxRepository(orderDB)

	cfg.DBName = "inventorydb"
	pgDSN = cfg.PostgresDSN()
	inventoryDB, err := pgxpool.New(ctx, pgDSN)
	// "postgres://orderuser:orderpass@localhost:5432/inventorydb"

	if err != nil {
		log.Fatal(err)
	}
	inventoryRepo := repository.NewPostgresOutboxRepository(inventoryDB)

	redisClient := redis.NewRedisClient(cfg.RedisAddr)

	orderPublisher := publisher.New(orderRepo, redisClient, "stream:orders")
	inventoryPublisher := publisher.New(inventoryRepo, redisClient, "stream:inventory")

	// outboxRepo := repository.NewPostgresOutboxRepository(pool)
	go func() {
		if err := orderPublisher.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := inventoryPublisher.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	select {}

}
