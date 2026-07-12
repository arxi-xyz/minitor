package socket

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	router := NewRouter()
	router.Handle("ping", func(client *Client, _ Message) {
		client.Send <- Message{Type: "pong"}
	})

	client := NewClient(conn, router)

	h.hub.Register(client)
	defer h.hub.Unregister(client.ID)

	defer func() {
		close(client.Send)
		conn.Close(websocket.StatusNormalClosure, "")
	}()

	ctx := r.Context()
	client.Run(ctx)
}
