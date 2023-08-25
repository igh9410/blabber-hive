package chat

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	hub *Hub
}

func NewWsHandler(h *Hub) *WsHandler {
	return &WsHandler{
		hub: h,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WsHandler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error occured : %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer func() {
		conn.Close()
	}()

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	log.Printf("chatroomId = ", roomID)

	if err != nil {
		log.Printf("Error occured with room ID %v: %v", roomID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	/*
		// Create a new client instance
		client := &Client{
			hub:  h.hub,
			conn: conn,
			send: make(chan []byte, 256),
			// Add chatroomID to the client if you have it in your Client struct
			chatroomID: roomID,
		}

		// Register the new client with the hub
		client.hub.register <- client

		// Start the read and write goroutines for the client
		go client.writePump()
		go client.readPump()
	*/
}
