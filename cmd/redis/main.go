package main

import (
	"context"
	"fmt"
	"order-system/internal/redis"
)

func main(){
	ctx := context.Background()

	redis := redis.NewRedisClient("localhost:6379")
	err := redis.Ping(ctx)

	fmt.Println(err)
}
