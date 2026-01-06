package auction

import "time"

// PlaceBid applies business rules for bidding
func (a *Auction) PlaceBid(amount int64, now time.Time) error {
	if a.Status == StatusClosed || now.After(a.EndTime) {
		a.Status = StatusClosed
		return ErrAuctionClosed
	}
	if amount <= a.HighestBid {
		return ErrBidTooLow
	}

	a.HighestBid = amount
	return nil
}

func (a *Auction) Close() {
	a.Status = StatusClosed
}
