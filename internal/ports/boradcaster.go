package ports

type Broadcaster interface {
	Broadcast(auctionID string, event any) error
}
