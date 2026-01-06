package ports

import (
	"auctioner/internal/domain/auction"
)

type AuctionRepository interface {
	Get(id string) (*auction.Auction, error)
	Save(a *auction.Auction) error
	ListOpen() ([]*auction.Auction, error)
}
