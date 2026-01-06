package application

import (
	"auctioner/internal/domain/auction"
	"testing"
	"time"
)

type mockRepo struct {
	auction *auction.Auction
}

func (m *mockRepo) Get(id string) (*auction.Auction, error) {
	return m.auction, nil
}

func (m *mockRepo) Save(a *auction.Auction) error {
	m.auction = a
	return nil
}

func (m *mockRepo) ListOpen() ([]*auction.Auction, error) {
	if m.auction != nil && m.auction.Status == auction.StatusOpen {
		return []*auction.Auction{m.auction}, nil
	}
	return []*auction.Auction{}, nil
}

type mockBroadcaster struct {
	events []map[string]any
}

func (m *mockBroadcaster) Broadcast(_ string, event any) error {
	m.events = append(m.events, event.(map[string]any))
	return nil
}

func TestPlaceBid_Success(t *testing.T) {
	repo := &mockRepo{
		auction: &auction.Auction{
			ID:         "a1",
			HighestBid: 100,
			Status:     auction.StatusOpen,
			EndTime:    time.Now().Add(time.Minute),
		},
	}

	bc := &mockBroadcaster{}
	svc := NewAuctionService(repo, bc)

	// freeze time
	svc.clock = func() time.Time {
		return time.Now()
	}

	err := svc.PlaceBid("a1", 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.auction.HighestBid != 200 {
		t.Fatalf("expected highest bid 200, got %d", repo.auction.HighestBid)
	}

	if len(bc.events) != 1 {
		t.Fatalf("expected 1 broadcast event")
	}
}

func TestPlaceBid_TooLow(t *testing.T) {
	repo := &mockRepo{
		auction: &auction.Auction{
			ID:         "a1",
			HighestBid: 100,
			Status:     auction.StatusOpen,
			EndTime:    time.Now().Add(time.Minute),
		},
	}

	bc := &mockBroadcaster{}
	svc := NewAuctionService(repo, bc)

	err := svc.PlaceBid("a1", 50)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if len(bc.events) != 0 {
		t.Fatalf("expected no broadcast")
	}
}

func TestPlaceBid_AuctionClosed(t *testing.T) {
	repo := &mockRepo{
		auction: &auction.Auction{
			ID:         "a1",
			HighestBid: 100,
			Status:     auction.StatusClosed,
			EndTime:    time.Now().Add(time.Minute),
		},
	}

	bc := &mockBroadcaster{}
	svc := NewAuctionService(repo, bc)

	err := svc.PlaceBid("a1", 200)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
