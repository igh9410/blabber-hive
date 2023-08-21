package chat

import (
	"time"
)

type Hub struct {
	Rooms      map[string]*ChatRoom
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*ChatRoom),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, roomExists := h.Rooms[client.ChatRoomID.String()]; !roomExists {

				h.Rooms[client.ChatRoomID.String()] = &ChatRoom{
					ID:        client.ChatRoomID,
					Clients:   make(map[string]*Client),
					CreatedAt: time.Now(),
				}
			}
			h.Rooms[client.ChatRoomID.String()].Clients[client.User.ID.String()] = client

		case client := <-h.Unregister:
			if room, roomExists := h.Rooms[client.ChatRoomID.String()]; roomExists {
				if _, clientExists := room.Clients[client.User.ID.String()]; clientExists {
					delete(room.Clients, client.User.ID.String())
				}

				// If the room is empty, you might consider deleting the room as well
				if len(room.Clients) == 0 {
					delete(h.Rooms, client.ChatRoomID.String())
				}
			}

		case msg := <-h.Broadcast:
			if room, roomExists := h.Rooms[msg.ChatRoomID.Variant().String()]; roomExists {
				for _, client := range room.Clients {
					select {
					case client.Send <- msg:
					default:
						// If there's an error sending the message to a client, you might consider removing the client
						delete(room.Clients, client.User.ID.String())
						close(client.Send)
					}
				}
			}
		}
	}
}
