package redis

import (
	"context"
	"encoding/json"
	"order-system/internal/domain"
	"strings"
	"time"

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

func (c *Client) AddToStream(
	ctx context.Context,
	stream string,
	values map[string]any,
) (string, error) {
	cmd := c.client.XAdd(ctx, &goredis.XAddArgs{
		Stream: stream,
		ID:     "*",
		Values: values,
	})
	return cmd.Result()
}

func (c *Client) Ack(
	ctx context.Context,
	stream string,
	group string,
	messageID string,

) error {
	return c.client.XAck(
		ctx,
		stream,
		group,
		messageID,
	).Err()
}

func (c *Client) Consume(
	ctx context.Context,
	stream string,
	group string,
	consumer string,
) ([]goredis.XStream, error) {
	cmd := c.client.XReadGroup(ctx,
		&goredis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams: []string{
				stream,
				">",
			},
			Count: 1,
			Block: 0 * time.Second,
		},
	)
	return cmd.Result()
}

func (c *Client) Decode(payload string) (domain.Order, error) {
	var order domain.Order
	err := json.Unmarshal([]byte(payload), &order)

	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (c *Client) EnsureGroup(
	ctx context.Context,
	stream string,
	group string,
) error {
	err := c.client.XGroupCreateMkStream(
		ctx,
		stream,
		group,
		"0",
	).Err()
	if err == nil {
		return nil
	}
	// Ignore "group already exists"
	if strings.Contains(err.Error(), "BUSYGROUP") {
		return nil
	}
	return err
}
