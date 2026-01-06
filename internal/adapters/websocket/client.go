package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	auctionID string
	hub       *Hub
}

func (c *Client) readPump(handler func([]byte)) {
	defer func() {
		c.hub.Unregister(c.auctionID, c)
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		handler(msg)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
