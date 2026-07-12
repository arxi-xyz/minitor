package socket

import (
	"context"
	"log"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan Message
	router *Router
}

func NewClient(conn *websocket.Conn, router *Router) *Client {
	return &Client{
		ID:     uuid.NewString(),
		Conn:   conn,
		Send:   make(chan Message),
		router: router,
	}
}

func (c *Client) WritingLoop(ctx context.Context) {
	for msg := range c.Send {
		data, err := msg.Marshal()
		if err != nil {
			log.Println(err)
			continue
		}

		if err := c.Conn.Write(ctx, websocket.MessageText, data); err != nil {
			log.Println(err)
			return
		}
	}
}

func (c *Client) ReadingLoop(ctx context.Context) {
	for {
		_, data, err := c.Conn.Read(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		msg, err := UnmarshalMessage(data)
		if err != nil {
			log.Println(err)
			continue
		}

		if msg.Type == "" {
			continue
		}

		c.router.Dispatch(c, msg)
	}
}

func (c *Client) Run(ctx context.Context) {
	go c.WritingLoop(ctx)
	c.ReadingLoop(ctx)
}
