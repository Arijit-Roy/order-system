package main

import (
	"context"
	"log"
	"order-system/internal/config"
	"order-system/internal/inventory"
	"order-system/internal/redis"
	redistransport "order-system/internal/transport/redis"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	cfg := config.Load()
	redisClient := redis.NewRedisClient(cfg.RedisAddr)
	ctx := context.Background()

	if err := redisClient.EnsureGroup(ctx, "stream:orders", "inventory-group"); err != nil {
		log.Fatal(err)
	}

	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	pgDSN := cfg.PostgresDSN()

	pool, err := pgxpool.New(ctx, pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := inventory.NewPostgresRepo(pool)

	inventoryService := inventory.NewService(repo)
	consumer := redistransport.NewConsumer(redisClient, "stream:orders", "inventory-group", "worder-1")
	if err := consumer.Run(ctx, inventoryService.Reserve); err != nil {
		log.Fatal(err)
	}

}
