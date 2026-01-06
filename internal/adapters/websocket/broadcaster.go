package websocket

import "encoding/json"

type Broadcaster struct {
	hub *Hub
}

func NewBroadcaster(hub *Hub) *Broadcaster {
	return &Broadcaster{hub: hub}
}

func (b *Broadcaster) Broadcast(auctionID string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	b.hub.Broadcast(auctionID, data)
	return nil
}
