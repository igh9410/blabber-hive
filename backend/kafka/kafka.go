package kafka

import (
	"backend/internal/chat"
	"encoding/json"

	"log"
	"os"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewKafkaClient() (*confluentKafka.AdminClient, error) {
	// Create a new AdminClient.
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	// Create a new AdminClient instance
	adminClient, err := confluentKafka.NewAdminClient(&confluentKafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Replace with the address of your Kafka broker
	})

	if err != nil {
		log.Printf("Failed to create Admin client: %s\n", err)
		return nil, err
	}
	/*
		// Contexts are used to abort or limit the amount of time
		// the Admin call blocks waiting for a result.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Define the topic configurations
		topic := "messages"

		numPartitions := 3
		replicationFactor := 1

		// Create topics on cluster.
		// Set Admin options to wait for the operation to finish (or at most 60s)
		maxDur, err := time.ParseDuration("60s")

		if err != nil {
			log.Printf("Failed to parse duration: %s\n", err)
		}

			// Check if the topic exists
			topicResults, err := adminClient.
			if err != nil {
				log.Printf("Failed to list topics: %v\n", err)
				return nil, err
			}

			results, err := adminClient.CreateTopics(
				ctx,
				// Multiple topics can be created simultaneously
				// by providing more TopicSpecification structs here.
				[]confluentKafka.TopicSpecification{{
					Topic:             topic,
					NumPartitions:     numPartitions,
					ReplicationFactor: replicationFactor}},
				// Admin options
				confluentKafka.SetAdminOperationTimeout(maxDur))

			if err != nil {
				log.Printf("Failed to create topic: %v\n", err)
				return nil, err
			}

			// Print results
			for _, result := range results {
				log.Printf("%s\n", result)
			} */

	return adminClient, nil
}

func KafkaProducer() (*confluentKafka.Producer, error) {
	producer, err := confluentKafka.NewProducer(&confluentKafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER_URL"), // Replace with address of your Kafka broker
	})

	if err != nil {
		log.Printf("Failed to create Kafka producer: %s", err)
		return nil, err
	}

	return producer, nil

}

func KafkaConsumer(batchProcessor *BatchProcessor) (*confluentKafka.Consumer, error) {
	consumer, err := confluentKafka.NewConsumer(&confluentKafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER_URL"),
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
