package auction

import "time"

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusClosed Status = "CLOSED"
)

type Auction struct {
	ID         string
	HighestBid int64
	Status     Status
	EndTime    time.Time
}
