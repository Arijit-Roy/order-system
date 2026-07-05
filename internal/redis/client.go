package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	client *goredis.Client
}

func NewRedisClient(addr string) *Client {
	return &Client{
		client: goredis.NewClient(
			&goredis.Options{
				Addr: addr,
			},
		),
	}
}

func (c *Client) Ping(
	ctx context.Context,
) error {
	return c.client.Ping(ctx).Err()
}