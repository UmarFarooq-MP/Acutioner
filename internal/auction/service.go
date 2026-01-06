package auction

import "time"

// PlaceBid applies business rules for bidding
func (a *Auction) PlaceBid(amount int64, now time.Time) error {
	// the current auction must be open
	if a.Status == StatusClosed || now.After(a.EndTime) {
		a.Status = StatusClosed
		return ErrAuctionClosed
	}

	// the current bid must be higher than the last succ
	if amount <= a.HighestBid {
		return ErrBidTooLow
	}

	// update bid
	a.HighestBid = amount
	return nil
}

func (a *Auction) Close() {
	a.Status = StatusClosed
}
