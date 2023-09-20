package chat

import (
	"log"
	"net/http"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	hub *Hub
	Service
	kafkaProducer *confluentKafka.Producer
}

func NewWsHandler(h *Hub, s Service, kafkaProducer *confluentKafka.Producer) *WsHandler {
	return &WsHandler{
		hub:           h,
		Service:       s,
		kafkaProducer: kafkaProducer,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WsHandler) RegisterClient(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error occured : %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	log.Printf("chatroomId = %s", roomID)

	if err != nil {
		log.Printf("Error occured with room ID %v: %v", roomID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Get the user ID from the context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID does not exist in the context"})
		return
	}

	// Validate the UUID format
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Create a new client instance
	// Register the new client
	client, err := h.Service.RegisterClient(c.Request.Context(), h.hub, conn, roomID, userID, h.kafkaProducer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occured with registering the client " + err.Error()})
		return
	}

	// Start the client's read and write pumps
	go client.writePump()
	go client.readPump()

}
