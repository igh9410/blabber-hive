package kafka

import (
	"backend/internal/chat"
	"encoding/json"

	"log"
	"os"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	bootstrapServersVar = "KAFKA_BROKER_URL"
)

func NewKafkaClient() (*confluentKafka.AdminClient, error) {
	// Create a new AdminClient.
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	// Create a new AdminClient instance
	bootstrapServers := os.Getenv(bootstrapServersVar)
	var config confluentKafka.ConfigMap

	log.Printf("Kafka bootstrap servers: %s", bootstrapServers)

	if os.Getenv("IS_PRODUCTION") == "YES" {
		config = confluentKafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
		}
		config["security.protocol"] = "SASL_SSL"
		config["sasl.mechanisms"] = "SCRAM-SHA-256"
		config["sasl.username"] = os.Getenv("KAFKA_USERNAME")
		config["sasl.password"] = os.Getenv("KAFKA_PASSWORD")
		//	config["debug"] = "all"

	} else {
		config = confluentKafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
		}
	}

	adminClient, err := confluentKafka.NewAdminClient(&config)
	if err != nil {
		log.Printf("Failed to create Admin client: %s\n", err)
		return nil, err
	}

	return adminClient, nil
}

func KafkaProducer() (*confluentKafka.Producer, error) {
	// Create a new Producer instance
	bootstrapServers := os.Getenv(bootstrapServersVar)
	var config confluentKafka.ConfigMap

	if os.Getenv("IS_PRODUCTION") == "YES" {
		config = confluentKafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
		}
		config["security.protocol"] = "SASL_SSL"
		config["sasl.mechanisms"] = "SCRAM-SHA-256"
		config["sasl.username"] = os.Getenv("KAFKA_USERNAME")
		config["sasl.password"] = os.Getenv("KAFKA_PASSWORD")
		//	config["debug"] = "all"

	} else {
		config = confluentKafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
		}
	}

	producer, err := confluentKafka.NewProducer(&config)
	if err != nil {
		log.Printf("Failed to create Kafka producer: %s", err)
		return nil, err
	}

	return producer, nil

}

func KafkaConsumer(batchProcessor *BatchProcessor) (*confluentKafka.Consumer, error) {
	// Create a new Consumer instance
	bootstrapServers := os.Getenv(bootstrapServersVar)
	var config confluentKafka.ConfigMap

	if os.Getenv("IS_PRODUCTION") == "YES" {
		config = confluentKafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
			"group.id":          "foo",
			"auto.offset.reset": "earliest",
		}
		config["security.protocol"] = "SASL_SSL"
		config["sasl.mechanisms"] = "SCRAM-SHA-256"
		config["sasl.username"] = os.Getenv("KAFKA_USERNAME")
		config["sasl.password"] = os.Getenv("KAFKA_PASSWORD")
		//	config["debug"] = "all"

	} else {
		config = confluentKafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
			"group.id":          "foo",
			"auto.offset.reset": "earliest",
		}
	}

	consumer, err := confluentKafka.NewConsumer(&config)
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
