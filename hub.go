package websocket

import "encoding/json"
import "bytes"

type EventHandler func(interface{})

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[*Client]bool
	onMessage  chan []byte
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	listeners  map[string][]EventHandler
}

// NewHub returns a new instance of Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		onMessage:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		listeners:  make(map[string][]EventHandler),
	}
}

// On adds an event listener for a specific event
func (h *Hub) On(event string, fn EventHandler) *Hub {
	if h.listeners[event] == nil {
		h.listeners[event] = []EventHandler{}
	}

	h.listeners[event] = append(h.listeners[event], fn)

	return h
}

type Msg struct {
	Type string
	Data interface{}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.onMessage:
			var msg Msg
			json.NewDecoder(bytes.NewReader(message)).Decode(msg)
			var fns = h.listeners[msg.Type]

			for _, fn := range fns {
				fn(msg.Data)
			}
		}
	}
}
