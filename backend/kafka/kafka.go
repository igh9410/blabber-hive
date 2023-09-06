package kafka

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewKafkaClient() (*kafka.AdminClient, error) {
	// Create a new AdminClient.
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	// Create a new AdminClient instance
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Replace with the address of your Kafka broker
	})

	if err != nil {
		log.Printf("Failed to create Admin client: %s\n", err)
		return nil, err
	}

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

	results, err := adminClient.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))

	if err != nil {
		log.Printf("Failed to create topic: %v\n", err)
		return nil, err
	}

	// Print results
	for _, result := range results {
		log.Printf("%s\n", result)
	}

	return adminClient, nil
}
