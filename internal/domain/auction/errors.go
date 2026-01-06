package auction

import "errors"

var (
	ErrAuctionClosed = errors.New("auction is closed")
	ErrBidTooLow     = errors.New("bid amount is too low")
)
