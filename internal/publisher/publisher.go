package publisher

import (
	"context"
	"log"
	"order-system/internal/redis"
	"time"
)

type Publisher struct {
	repo   Repository
	redis  *redis.Client
	stream string
}

func New(repo Repository, redis *redis.Client, stream string) *Publisher {
	return &Publisher{
		repo:   repo,
		redis:  redis,
		stream: stream,
	}
}

func (p *Publisher) Run(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("outbox publisher started")

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := p.publishPending(ctx); err != nil {
				log.Println(err)
			}
		}

	}
}

func (p *Publisher) publishPending(ctx context.Context) error {
	events, err := p.repo.ListPending(ctx, 2000)

	if err != nil {
		return err
	}

	for _, event := range events {

		_, err := p.redis.AddToStream(
			ctx,
			p.stream,
			map[string]any{
				"event_type": event.EventType,
				"payload":    string(event.Payload),
			},
		)
		if err != nil {
			log.Println(err)
			continue
		}
		if err := p.repo.MarkPublished(
			ctx,
			event.ID,
		); err != nil {
			log.Println(err)
		}
	}
	return nil
}
