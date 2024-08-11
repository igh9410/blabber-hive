package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/igh9410/blabber-hive/backend/db"
	"github.com/igh9410/blabber-hive/backend/infra/kafka"
	"github.com/igh9410/blabber-hive/backend/infra/redis"
	"github.com/igh9410/blabber-hive/backend/internal/chat"
	"github.com/igh9410/blabber-hive/backend/internal/match"
	"github.com/igh9410/blabber-hive/backend/internal/user"
	"github.com/igh9410/blabber-hive/backend/router"

	"time"

	"github.com/joho/godotenv"
)

// @title           Blabber-Hive API
// @version         1.0
// @description     Blabber-Hive 문서
// @termsOfService  http://swagger.io/terms/
// @contact.name   임건혁
// @contact.url    http://www.swagger.io/support
// @contact.email  athanasia9410@gmail.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /api
// @securityDefinitions.basic  BasicAuth
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	if err := godotenv.Load(".env"); err != nil { // Running in local, must be run on go run . in ./cmd directory
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
		log.Fatalf("Failed to initialize Kafka cluster connection: %s", err)
	}
	defer kafkaClient.Close()

	kafkaProducer, err := kafka.KafkaProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %s", err)
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
		log.Fatalf("Failed to initialize Kafka consumer: %s", err)
	}

	routerConfig := &router.RouterConfig{
		UserHandler:   userHandler,
		ChatHandler:   chatHandler,
		ChatWsHandler: chatWsHandler,
		MatchHandler:  matchHandler,
		// Future handlers can be added here without changing the InitRouter signature
	}

	r := router.InitRouter(routerConfig)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server listening on http://localhost:8080")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("Server exiting")
}
