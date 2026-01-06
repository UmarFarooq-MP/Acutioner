package ports

import "auctioner/internal/auction"

type AuctionRepository interface {
	Get(id string) (*auction.Auction, error)
	Save(a *auction.Auction) error
}
