package main

import (
	"backend/db"
	"backend/internal/chat"
	"backend/internal/user"
	"backend/kafka"
	"backend/router"
	"log"

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

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	chatRep := chat.NewRepository(dbConn.GetDB())
	chatSvc := chat.NewService(chatRep, userRep)
	chatHandler := chat.NewHandler(chatSvc, userSvc)

	hub := chat.NewHub()

	go hub.Run()

	chatWsHandler := chat.NewWsHandler(hub, chatSvc)

	routerConfig := &router.RouterConfig{
		UserHandler:   userHandler,
		ChatHandler:   chatHandler,
		ChatWsHandler: chatWsHandler,
		// Future handlers can be added here without changing the InitRouter signature
	}

	router.InitRouter(routerConfig)

}
