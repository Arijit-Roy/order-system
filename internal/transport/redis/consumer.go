package redistransport

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"order-system/internal/domain"
	"order-system/internal/redis"

	goredis "github.com/redis/go-redis/v9"
)

type Handler func(ctx context.Context, order domain.Order) error

type Consumer struct {
	redis    *redis.Client
	stream   string
	group    string
	consumer string
}

func NewConsumer(
	redis *redis.Client,
	stream string,
	group string,
	consumer string,
) *Consumer {
	return &Consumer{
		redis:    redis,
		stream:   stream,
		group:    group,
		consumer: consumer,
	}
}

func (c *Consumer) Run(ctx context.Context, handler Handler) error {
	for {
		streams, err := c.redis.Consume(ctx, c.stream, c.group, c.consumer)

		if err != nil {
			return err
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				order, err := decodeMessage(message)
				// c.redis.Decode(message)

				if err != nil {
					log.Println(err)
					continue
				}
				if err := handler(ctx, order); err != nil {
					log.Println(err)
					continue
				}
				if err := c.redis.Ack(
					ctx,
					c.stream,
					c.group,
					message.ID,
				); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func decodeMessage(message goredis.XMessage) (domain.Order, error) {
	payload, ok := message.Values["payload"].(string)

	if !ok {
		return domain.Order{}, errors.New("payload missing")
	}

	var order domain.Order

	err := json.Unmarshal([]byte(payload), &order)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
