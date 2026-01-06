package cmd

import (
	"context"
	"log"
	"net/http"
	"time"

	redisAdapter "auctioner/internal/adapters/redis"
	"auctioner/internal/adapters/repository"
	ws "auctioner/internal/adapters/websocket"
	"auctioner/internal/application"
	"auctioner/internal/domain/auction"

	"github.com/redis/go-redis/v9"
)

func startRedisServer(port string) {
	log.Println("Starting auction server (redis mode)")

	repo := repository.NewInMemoryAuctionRepository()

	// Adapter-only seeding (demo / dev)
	memRepo := repo.(*repository.InMemoryAuctionRepository)
	memRepo.Seed(&auction.Auction{
		ID:         "auction-1",
		HighestBid: 100,
		Status:     auction.StatusOpen,
		EndTime:    time.Now().Add(1 * time.Minute),
	})

	hub := ws.NewHub()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Docker Compose service name
	})

	// Fail fast if Redis is not available
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis not reachable: %v", err)
	}

	log.Println("Connected to Redis")

	broadcaster := redisAdapter.NewPublisher(redisClient)

	subscriber := redisAdapter.NewSubscriber(redisClient, hub)
	subscriber.Start(context.Background())

	auctionService := application.NewAuctionService(repo, broadcaster)

	wsHandler := ws.NewHandler(hub, auctionService)

	// Worker to remove time out listings
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		log.Println("⏱ Timeout worker started")

		for range ticker.C {
			if err := auctionService.CloseExpiredAuctions(); err != nil {
				log.Println("timeout worker error:", err)
			}
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/ws", wsHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("WebSocket endpoint available at /ws")
	log.Println("Listening on port", port)
	log.Fatal(server.ListenAndServe())
}
