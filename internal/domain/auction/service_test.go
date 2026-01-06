package auction

import (
	"errors"
	"testing"
	"time"
)

func TestPlaceBid_Success(t *testing.T) {
	a := &Auction{
		ID:         "1",
		HighestBid: 100,
		Status:     StatusOpen,
		EndTime:    time.Now().Add(time.Minute),
	}

	err := a.PlaceBid(200, time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if a.HighestBid != 200 {
		t.Fatalf("expected highest bid 200, got %d", a.HighestBid)
	}
}

func TestPlaceBid_TooLow(t *testing.T) {
	a := &Auction{
		HighestBid: 100,
		Status:     StatusOpen,
		EndTime:    time.Now().Add(time.Minute),
	}

	err := a.PlaceBid(50, time.Now())
	if !errors.Is(err, ErrBidTooLow) {
		t.Fatalf("expected ErrBidTooLow, got %v", err)
	}
}

func TestPlaceBid_AuctionClosed(t *testing.T) {
	a := &Auction{
		HighestBid: 100,
		Status:     StatusClosed,
		EndTime:    time.Now().Add(time.Minute),
	}

	err := a.PlaceBid(200, time.Now())
	if !errors.Is(err, ErrAuctionClosed) {
		t.Fatalf("expected ErrAuctionClosed, got %v", err)
	}
}
