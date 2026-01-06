package repository

import (
	"auctioner/internal/domain/auction"
	"auctioner/internal/ports"
	"errors"
	"sync"
)

var ErrAuctionNotFound = errors.New("auction not found")

type InMemoryAuctionRepository struct {
	mu       sync.Mutex
	auctions map[string]*auction.Auction
}

func NewInMemoryAuctionRepository() ports.AuctionRepository {
	return &InMemoryAuctionRepository{
		auctions: make(map[string]*auction.Auction),
	}
}

func (r *InMemoryAuctionRepository) Get(id string) (*auction.Auction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	a, ok := r.auctions[id]
	if !ok {
		return nil, ErrAuctionNotFound
	}
	return a, nil
}

func (r *InMemoryAuctionRepository) Save(a *auction.Auction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.auctions[a.ID] = a
	return nil
}

// Seed adapter-only helper (demo / tests)
func (r *InMemoryAuctionRepository) Seed(a *auction.Auction) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.auctions[a.ID] = a
}

func (r *InMemoryAuctionRepository) ListOpen() ([]*auction.Auction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var res []*auction.Auction
	for _, a := range r.auctions {
		if a.Status == auction.StatusOpen {
			res = append(res, a)
		}
	}
	return res, nil
}
