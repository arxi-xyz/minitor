package socket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Handler struct {
	hub     *Hub
	monitor *Monitor
	router  *Router
}

func NewHandler(hub *Hub, monitor *Monitor) *Handler {
	h := &Handler{
		hub:     hub,
		monitor: monitor,
		router:  NewRouter(),
	}

	h.router.Handle("ping", func(client *Client, _ Message) {
		client.Send <- Message{Type: "pong"}
	})

	h.router.Handle("processes", func(client *Client, msg Message) {
		var req ProcessesRequest
		if len(msg.Data) > 0 {
			if err := json.Unmarshal(msg.Data, &req); err != nil {
				return
			}
		}

		page := h.monitor.ProcessesPage(req.Offset, req.Limit)
		resp, err := page.Message()
		if err != nil {
			log.Println(err)
			return
		}

		client.Send <- resp
	})

	return h
}

func (h *Handler) Router() *Router {
	return h.router
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, h.router)

	h.hub.Register(client)
	defer h.hub.Unregister(client.ID)

	defer func() {
		close(client.Send)
		conn.Close(websocket.StatusNormalClosure, "")
	}()

	if msg, ok := h.hub.LastMessage(); ok {
		select {
		case client.Send <- msg:
		default:
		}
	}

	ctx := r.Context()
	client.Run(ctx)
}
