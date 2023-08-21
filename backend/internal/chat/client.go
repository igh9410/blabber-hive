package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error writing to client %v: %v", c.User.Username, err)
			break
		}
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:   string(m),
			CreatedAt: time.Now(),
		}

		hub.Broadcast <- msg
	}
}
