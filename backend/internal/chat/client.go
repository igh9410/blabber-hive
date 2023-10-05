package chat

import (
	"encoding/json"
	"log"
	"time"
	"unicode/utf8"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub           *Hub
	conn          *websocket.Conn
	send          chan []byte
	chatroomID    uuid.UUID
	senderID      uuid.UUID
	kafkaProducer *confluentKafka.Producer
}

const MaxMessageSize = 1000
const ErrorMessage = "Message is too long. Can't exceed 1000 characters."

func NewClient(hub *Hub, conn *websocket.Conn, chatroomID uuid.UUID, senderID uuid.UUID, kafkaProducer *confluentKafka.Producer) *Client {
	return &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, 256),
		chatroomID:    chatroomID,
		senderID:      senderID,
		kafkaProducer: kafkaProducer,
	}
}

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

		messageContent := <-c.send
		// Check message length in terms of characters
		if utf8.RuneCountInString(string(messageContent)) > MaxMessageSize {
			log.Printf("Message too long: %d characters", utf8.RuneCountInString(string(messageContent)))

			// Send an error message to the client
			c.send <- []byte(ErrorMessage)

			continue
		}

		log.Printf("Write message = %v", string(messageContent))
		err := c.conn.WriteMessage(websocket.TextMessage, messageContent)
		if err != nil {
			log.Println("write:", err)
			break
		}

		// Construct the message
		msg := &Message{
			ChatRoomID: c.chatroomID,           // Using the client's chatroomID
			SenderID:   c.senderID,             // Using the client's senderID
			Content:    string(messageContent), // Using the received message content
			CreatedAt:  time.Now(),             // Using the current time for CreatedAt
			// Other fields can be set as per your requirements
		}

		// Convert the message to JSON
		jsonMessage, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)

			continue
		}

		// Send the message to the Kafka topic "messages"
		topic := "messages"
		kafkaMessage := &confluentKafka.Message{
			TopicPartition: confluentKafka.TopicPartition{
				Topic:     &topic, // Use the Kafka topic name "messages"
				Partition: confluentKafka.PartitionAny,
			},
			Value: []byte(jsonMessage),
		}
		err = c.kafkaProducer.Produce(kafkaMessage, nil)

		if err != nil {
			log.Printf("Failed to produce message to Kafka: %v", err)
		} else {
			log.Printf("Produced message to Kafka: %v", kafkaMessage.Value)
		}
	}
}
