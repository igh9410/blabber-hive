package chat

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	chatroomID uuid.UUID
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		log.Printf("Read message = %v", message)
		if err != nil {
			log.Println("read:", err)
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	for {
		message := <-c.send
		log.Printf("Read message = %v", message)
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
