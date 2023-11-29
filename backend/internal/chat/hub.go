package chat

import (
	"log"

	"github.com/google/uuid"
)

type BroadcastMessage struct {
	Message []byte
	Sender  *Client
}

type Hub struct {
	clients map[uuid.UUID][]*Client

	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uuid.UUID][]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.chatroomID] = append(h.clients[client.chatroomID], client)
			log.Println("Client roomID:", client.chatroomID)
			log.Println("Client registered:", client)

		case client := <-h.unregister:
			if clientsInRoom, ok := h.clients[client.chatroomID]; ok {
				for i, c := range clientsInRoom {
					if c == client {
						// Remove client from slice
						h.clients[client.chatroomID] = append(clientsInRoom[:i], clientsInRoom[i+1:]...)
						close(client.send)
						log.Println("Client unregistered:", client)
						break
					}
				}
				if len(h.clients[client.chatroomID]) == 0 {
					delete(h.clients, client.chatroomID)
				}
			}

		case broadcastMsg := <-h.broadcast:
			if clientsInRoom, ok := h.clients[broadcastMsg.Sender.chatroomID]; ok {
				log.Printf("Broadcasting message to clients in chat room: %s", broadcastMsg.Sender.chatroomID)
				for _, client := range clientsInRoom {
					if client != broadcastMsg.Sender {
						select {
						case client.send <- broadcastMsg.Message:
							// Message sent successfully
						default:
							// Handle failed message send, e.g., close connection and remove client
							close(client.send)
							// Remove client from slice
							// Note: This requires careful handling to avoid modifying the slice while iterating over it
						}
					}
				}
			}
		}
	}
}
