package main

import (
	"context"
	"fmt"
	"order-system/internal/config"
	"order-system/internal/redis"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	redis := redis.NewRedisClient(cfg.RedisAddr)
	err := redis.Ping(ctx)

	fmt.Println(err)
}
