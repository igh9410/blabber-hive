package main

import (
	"backend/db"
	"backend/internal/chat"
	"backend/internal/user"
	"backend/kafka"
	"backend/router"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil { // Running in local, must be run on go run . in ./cmd directory
		log.Println("No .env file found. Using OS environment variables.")
	}

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing the database: %s", err)
		}
	}()
	log.Println("Database initialized")

	kafkaClient, err := kafka.NewKafkaClient()
	if err != nil {
		log.Printf("Failed to initialize Kafka cluster connection")
	}
	defer kafkaClient.Close()

	kafkaProducer, err := kafka.KafkaProducer()
	if err != nil {
		log.Printf("Failed to initialize Kafka producer")
	}
	defer kafkaProducer.Close()

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	chatRep := chat.NewRepository(dbConn.GetDB())
	chatSvc := chat.NewService(chatRep, userRep, kafkaProducer)
	chatHandler := chat.NewHandler(chatSvc, userSvc)

	hub := chat.NewHub()

	go hub.Run()

	chatWsHandler := chat.NewWsHandler(hub, chatSvc, kafkaProducer)

	// Initialize BatchProcessor with a function to insert messages into Postgres and start KafkaConsumer
	batchProcessor := kafka.NewBatchProcessor(func(messages []chat.Message) error {
		// Code to insert messages into Postgres
		// Ensure this function is safe to be called concurrently
		return nil
	}, 100, 5*time.Second)

	if _, err := kafka.KafkaConsumer(batchProcessor); err != nil {
		log.Printf("Failed to initialize Kafka consumer: %s", err)
	}

	routerConfig := &router.RouterConfig{
		UserHandler:   userHandler,
		ChatHandler:   chatHandler,
		ChatWsHandler: chatWsHandler,
		// Future handlers can be added here without changing the InitRouter signature
	}

	router.InitRouter(routerConfig)

}
