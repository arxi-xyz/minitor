package socket

import "sync"

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client

	lastMu sync.RWMutex
	last   Message
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.ID] = client
}

func (h *Hub) Unregister(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, id)
}

func (h *Hub) Get(id string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, ok := h.clients[id]
	return client, ok
}

func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.clients)
}

func (h *Hub) Broadcast(msg Message) {
	h.lastMu.Lock()
	h.last = msg
	h.lastMu.Unlock()

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.Send <- msg:
		default:
		}
	}
}

func (h *Hub) LastMessage() (Message, bool) {
	h.lastMu.RLock()
	defer h.lastMu.RUnlock()

	if h.last.Type == "" {
		return Message{}, false
	}

	return h.last, true
}

func (h *Hub) Clients() []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := make([]*Client, 0, len(h.clients))

	for _, client := range h.clients {
		clients = append(clients, client)
	}

	return clients
}
