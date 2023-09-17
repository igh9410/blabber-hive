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

	// Retrieve email from the context
	senderEmail, exists := c.Get("email")
	if !exists {
		log.Printf("Email not found in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email or Email Not Found"})
		return
	}

	emailStr, ok := senderEmail.(string)
	if !ok {
		log.Println("Failed to assert userEmail as string")
		return
	}

	// Create a new client instance
	// Register the new client
	client, err := h.Service.RegisterClient(c.Request.Context(), h.hub, conn, roomID, emailStr, h.kafkaProducer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occured with registering the client " + err.Error()})
		return
	}

	go client.writePump()
	go client.readPump()

}
