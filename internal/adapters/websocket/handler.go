package websocket

import (
	"auctioner/internal/application"
	"auctioner/internal/domain/auction"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type incomingMessage struct {
	Type      string `json:"type"`
	AuctionID string `json:"auction_id"`
	Amount    int64  `json:"amount"`
}

type Handler struct {
	upgrader websocket.Upgrader
	hub      *Hub
	service  *application.AuctionService
}

func NewHandler(hub *Hub, svc *application.AuctionService) *Handler {
	return &Handler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		hub:     hub,
		service: svc,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auctionID := r.URL.Query().Get("auction_id")
	if auctionID == "" {
		http.Error(w, "missing auction_id", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		conn:      conn,
		send:      make(chan []byte, 16),
		auctionID: auctionID,
		hub:       h.hub,
	}

	h.hub.Register(auctionID, client)

	go client.writePump()
	go client.readPump(func(msg []byte) {
		h.handleMessage(client, msg)
	})
}

func (h *Handler) handleMessage(c *Client, raw []byte) {
	var msg incomingMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}

	switch msg.Type {
	case "PLACE_BID":
		if err := h.service.PlaceBid(msg.AuctionID, msg.Amount); err != nil {
			var reason string
			switch err {
			case auction.ErrAuctionClosed:
				reason = "AUCTION_CLOSED"
			case auction.ErrBidTooLow:
				reason = "BID_TOO_LOW"
			default:
				reason = "UNKNOWN_ERROR"
			}

			// send error ONLY to this client
			resp := map[string]any{
				"type":       "BID_REJECTED",
				"auction_id": msg.AuctionID,
				"reason":     reason,
			}
			data, _ := json.Marshal(resp)
			c.send <- data
			return
		}
	}
}
