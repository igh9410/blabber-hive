package main

import (
	"backend/db"
	"backend/infra/kafka"
	"backend/infra/redis"
	"backend/internal/chat"
	"backend/internal/match"
	"backend/internal/user"
	"backend/router"
	"context"
	"log"
	"log/slog"
	"os"
	"strconv"

	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil { // Running in local, must be run on go run . in ./cmd directory
		slog.Info("No .env file found. Using OS environment variables.")
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

	redisClient, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("Could not initialize Redis connection: %s", err)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing the Redis connection: %s", err)
		}
	}()

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

	matchRep := match.NewRepository(redisClient)
	matchSvc := match.NewService(matchRep, chatRep, redisClient)
	matchHandler := match.NewHandler(matchSvc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Implementing graceful shutdown of Redis client

	go matchSvc.StartMatchmakingSubscriber(ctx)

	// Create an insert function with the database connection
	insertFunc := kafka.NewInsertFunc(dbConn.GetDB())

	// Initialize BatchProcessor with a function to insert messages into Postgres and start KafkaConsumer
	bulkInsertSizeStr := os.Getenv("BULK_INSERT_SIZE")
	bulkInsertSize, err := strconv.Atoi(bulkInsertSizeStr)
	if err != nil {
		log.Printf("Invalid integer for BULK_INSERT_SIZE: %v", err)
	}

	bulkInsertTimeStr := os.Getenv("BULK_INSERT_TIME")
	bulkInsertTime, err := time.ParseDuration(bulkInsertTimeStr)
	if err != nil {
		log.Printf("Invalid duration for BULK_INSERT_TIME: %v", err)
	}

	batchProcessor := kafka.NewBatchProcessor(insertFunc, bulkInsertSize, bulkInsertTime)

	defer batchProcessor.Stop()

	if _, err := kafka.KafkaConsumer(batchProcessor); err != nil {
		log.Printf("Failed to initialize Kafka consumer: %s", err)
	}

	routerConfig := &router.RouterConfig{
		UserHandler:   userHandler,
		ChatHandler:   chatHandler,
		ChatWsHandler: chatWsHandler,
		MatchHandler:  matchHandler,
		// Future handlers can be added here without changing the InitRouter signature
	}

	router.InitRouter(routerConfig)
	cancel()
}
