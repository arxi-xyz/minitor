package socket

type Router struct {
	handlers map[string]HandlerF
}

type HandlerF func(*Client, Message)

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerF),
	}
}

func (r *Router) Handle(name string, handler HandlerF) {
	r.handlers[name] = handler
}

func (r *Router) Dispatch(client *Client, msg Message) {
	handler, ok := r.handlers[msg.Type]

	if !ok {
		return
	}

	handler(client, msg)
}
