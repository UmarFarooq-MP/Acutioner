package cmd

import (
	"auctioner/internal/adapters/repository"
	ws "auctioner/internal/adapters/websocket"
	"auctioner/internal/application"
	"auctioner/internal/domain/auction"
	"log"
	"net/http"
	"time"
)

func startInMemoryServer(port string) {
	log.Println("Starting server (in-memory mode)")

	repo := repository.NewInMemoryAuctionRepository()
	memRepo := repo.(*repository.InMemoryAuctionRepository)

	memRepo.Seed(&auction.Auction{
		ID:         "auction-1",
		HighestBid: 100,
		Status:     auction.StatusOpen,
		EndTime:    time.Now().Add(10 * time.Minute),
	})

	hub := ws.NewHub()
	broadcaster := ws.NewBroadcaster(hub)

	service := application.NewAuctionService(repo, broadcaster)
	handler := ws.NewHandler(hub, service)

	http.Handle("/ws", handler)

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
