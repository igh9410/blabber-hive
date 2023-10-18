package kafka

import (
	"backend/internal/chat"
	"encoding/json"
	"fmt"

	"log"
	"os"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	bootstrapServersVar = "KAFKA_BROKER_URL"
)

var bootstrapServers = os.Getenv(bootstrapServersVar)

func NewKafkaClient() (*confluentKafka.AdminClient, error) {
	// Create a new AdminClient.
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	// Create a new AdminClient instance

	adminClient, err := confluentKafka.NewAdminClient(&confluentKafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s,localhost:9092", bootstrapServers), // Replace with the address of your Kafka broker
	})

	if err != nil {
		log.Printf("Failed to create Admin client: %s\n", err)
		return nil, err
	}

	return adminClient, nil
}

func KafkaProducer() (*confluentKafka.Producer, error) {
	producer, err := confluentKafka.NewProducer(&confluentKafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s,localhost:9092", bootstrapServers), // Replace with address of your Kafka broker
	})

	if err != nil {
		log.Printf("Failed to create Kafka producer: %s", err)
		return nil, err
	}

	return producer, nil

}

func KafkaConsumer(batchProcessor *BatchProcessor) (*confluentKafka.Consumer, error) {
	consumer, err := confluentKafka.NewConsumer(&confluentKafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s,localhost:9092", bootstrapServers), // Replace with address of your Kafka broker
		"group.id":          "foo",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Printf("Failed to create Kafka consumer: %s", err)
		return nil, err
	}

	// Subscribe to topic
	err = consumer.Subscribe("messages", nil)
	if err != nil {
		log.Printf("Failed to subscribe to topic: %s\n", err)
		return nil, err
	}

	// Run a separate goroutine to poll for messages
	go func() {
		for {
			msg, err := consumer.ReadMessage(-1) // Block indefinitely
			if err != nil {
				log.Printf("Consumer error: %v (%v)\n", err, msg)
				continue
			}

			// Add the message to the batch processor
			var message chat.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				log.Printf("Failed to unmarshal JSON: %v", err)
				continue
			}

			batchProcessor.AddMessage(message)
		}
	}()

	return consumer, nil
}
