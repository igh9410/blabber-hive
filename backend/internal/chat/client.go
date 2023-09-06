package chat

import (
	"log"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	chatroomID uuid.UUID
	senderID   uuid.UUID
}

const MaxMessageSize = 1000
const ErrorMessage = "Message is too long. Can't exceed 1000 characters."

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			break
		}
		// Check message length in terms of characters
		if utf8.RuneCountInString(string(message)) > MaxMessageSize {
			log.Printf("Message too long: %d characters", utf8.RuneCountInString(string(message)))

			// Send an error message to the client
			c.send <- []byte(ErrorMessage)

			continue
		}

		log.Printf("Read message = %v", string(message))

		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	for {

		message := <-c.send
		// Check message length in terms of characters
		if utf8.RuneCountInString(string(message)) > MaxMessageSize {
			log.Printf("Message too long: %d characters", utf8.RuneCountInString(string(message)))

			// Send an error message to the client
			c.send <- []byte(ErrorMessage)

			continue
		}

		log.Printf("Write message = %v", string(message))
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
