package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	client *redis.Client
}

func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) Broadcast(auctionID string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	channel := "auction:" + auctionID
	return p.client.Publish(context.Background(), channel, data).Err()
}
