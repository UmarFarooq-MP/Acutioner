package websocket

import "sync"

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*Client]struct{}),
	}
}

func (h *Hub) Register(auctionID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[auctionID]; !ok {
		h.clients[auctionID] = make(map[*Client]struct{})
	}
	h.clients[auctionID][c] = struct{}{}
}

func (h *Hub) Unregister(auctionID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[auctionID]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.clients, auctionID)
		}
	}
}

func (h *Hub) Broadcast(auctionID string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients[auctionID] {
		select {
		case c.send <- msg:
		default:
			// TODO:: drop slow client
		}
	}
}
