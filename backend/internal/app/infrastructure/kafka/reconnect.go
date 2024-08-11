package kafka

import (
	"log"
	"time"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// ReconnectLoop continuously attempts to establish a Kafka connection
// and sends a signal on the returned channel when a connection is made.
// The loop will continue to run until the provided stop channel is closed.
// It also returns the Kafka producer instance.
func ReconnectLoop(stop chan struct{}) (chan struct{}, *confluentKafka.Producer) {
	connected := make(chan struct{})
	var producer *confluentKafka.Producer

	go func() {
		for {
			select {
			case <-stop:
				log.Println("Stopping Kafka reconnection loop")
				close(connected)
				if producer != nil {
					producer.Close()
				}
				return
			default:
				kafkaClient, err := NewKafkaClient()
				if err != nil {
					log.Printf("Failed to initialize Kafka cluster connection: %s", err)
					time.Sleep(30 * time.Second) // Wait for 30 seconds before retrying
					continue
				}
				defer kafkaClient.Close()

				producer, err = KafkaProducer()
				if err != nil {
					log.Printf("Failed to initialize Kafka producer: %s", err)
					time.Sleep(30 * time.Second) // Wait for 30 seconds before retrying
					continue
				}

				log.Println("Kafka connection established")
				connected <- struct{}{}

				// Wait for the connection to be closed (e.g., due to an error)
				<-connected
			}
		}
	}()

	return connected, producer
}
