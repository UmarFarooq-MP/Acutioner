package application

import (
	"auctioner/internal/ports"
	"time"
)

type AuctionService struct {
	repo        ports.AuctionRepository
	broadcaster ports.Broadcaster
	clock       func() time.Time
}

func NewAuctionService(
	repo ports.AuctionRepository,
	broadcaster ports.Broadcaster,
) *AuctionService {
	return &AuctionService{
		repo:        repo,
		broadcaster: broadcaster,
		clock:       time.Now,
	}
}

// PlaceBid executes the bidding use-case
func (s *AuctionService) PlaceBid(
	auctionID string,
	amount int64,
) error {

	// load auction
	a, err := s.repo.Get(auctionID)
	if err != nil {
		return err
	}

	// place bid via interface
	if err := a.PlaceBid(amount, s.clock()); err != nil {
		return err
	}

	// updated state
	if err := s.repo.Save(a); err != nil {
		return err
	}

	// broadcast event (best-effort)
	event := map[string]any{
		"type":       "NEW_BID",
		"auction_id": a.ID,
		"amount":     a.HighestBid,
	}
	_ = s.broadcaster.Broadcast(a.ID, event)
	return nil
}

// CloseAuction ends the auction and broadcasts the result
func (s *AuctionService) CloseAuction(auctionID string) error {
	a, err := s.repo.Get(auctionID)
	if err != nil {
		return err
	}

	a.Close()

	if err := s.repo.Save(a); err != nil {
		return err
	}

	event := map[string]any{
		"type":       "AUCTION_ENDED",
		"auction_id": a.ID,
		"final_bid":  a.HighestBid,
	}
	_ = s.broadcaster.Broadcast(a.ID, event)

	return nil
}

func (s *AuctionService) CloseExpiredAuctions() error {
	auctions, err := s.repo.ListOpen()
	if err != nil {
		return err
	}

	now := s.clock()
	for _, a := range auctions {
		if now.After(a.EndTime) {
			a.Close()
			if err := s.repo.Save(a); err != nil {
				return err
			}
			event := map[string]any{
				"type":       "AUCTION_TIMED_OUT",
				"auction_id": a.ID,
				"final_bid":  a.HighestBid,
			}

			// best-effort broadcast
			_ = s.broadcaster.Broadcast(a.ID, event)
		}
	}
	return nil
}
