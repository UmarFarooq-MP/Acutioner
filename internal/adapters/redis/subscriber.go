package redis

import (
	"context"
	"log"

	"auctioner/internal/adapters/websocket"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	client *redis.Client
	hub    *websocket.Hub
}

func NewSubscriber(client *redis.Client, hub *websocket.Hub) *Subscriber {
	return &Subscriber{
		client: client,
		hub:    hub,
	}
}

// Start Subscribe to all auction channels
func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.client.PSubscribe(ctx, "auction:*")

	go func() {
		ch := pubsub.Channel()

		for msg := range ch {
			auctionID := msg.Channel[len("auction:"):]
			s.hub.Broadcast(auctionID, []byte(msg.Payload))
		}
	}()

	log.Println("Redis subscriber started")
}
