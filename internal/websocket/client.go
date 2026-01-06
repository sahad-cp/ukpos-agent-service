package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

func (c *Client) Close() {
	panic("unimplemented")
}

func Connect(url string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &Client{Conn: conn}, nil
}

func (c *Client) Send(v any) {
	_ = c.Conn.WriteJSON(v)
}

func (c *Client) Listen(handler func(map[string]any)) {
	for {
		var msg map[string]any
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println("WS error:", err)
			return
		}
		handler(msg)
	}
}
