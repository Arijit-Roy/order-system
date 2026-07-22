package main

import (
	"context"
	"log"
	"order-system/internal/config"
	"order-system/internal/notification"
	"order-system/internal/redis"
	redistransport "order-system/internal/transport/redis"
)

func main() {
	cfg := config.Load()
	redisClient := redis.NewRedisClient(cfg.RedisAddr)
	ctx := context.Background()

	if err := redisClient.EnsureGroup(ctx, "stream:orders", "notification-group"); err != nil {
		log.Fatal(err)
	}

	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	notifService := notification.NewService()
	consumer := redistransport.NewConsumer(redisClient, "stream:orders", "notification-group", "worker-1")

	consumer.Run(ctx, notifService.SendOrderConfirmation)
}
