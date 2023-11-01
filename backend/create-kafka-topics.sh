#!/bin/bash

# Kafka broker address
KAFKA_BROKER="broker:29092"

# Create the 'messages' topic
/usr/bin/kafka-topics --create --bootstrap-server $KAFKA_BROKER --topic messages --partitions 3 --replication-factor 1
# Add more topics as needed
# /usr/bin/kafka-topics --create --bootstrap-server $KAFKA_BROKER --topic messages --partitions 3 --replication-factor 1
