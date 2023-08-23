package main

import (
	"backend/db"

	"backend/internal/chat"
	"backend/internal/user"
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

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	chatRep := chat.NewRepository(dbConn.GetDB())
	chatSvc := chat.NewService(chatRep)
	chatHandler := chat.NewHandler(chatSvc)

	hub := chat.NewHub()
	wsHandler := chat.NewWsHandler(hub)

	router.InitRouter(userHandler, chatHandler, wsHandler)

}
