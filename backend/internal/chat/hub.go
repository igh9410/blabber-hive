package chat

import "log"

type BroadcastMessage struct {
	Message []byte
	Sender  *Client
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Client registered:", client)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client unregistered:", client)
			}
		case broadcastMsg := <-h.broadcast:
			log.Printf("Broadcasting message to %d clients: %s", len(h.clients), broadcastMsg.Message)
			for client := range h.clients {
				// Skip the sender
				if client == broadcastMsg.Sender {
					continue
				}
				select {
				case client.send <- broadcastMsg.Message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
