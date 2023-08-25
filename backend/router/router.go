package router

import (
	"backend/api/middleware"

	"backend/internal/chat"
	"backend/internal/user"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

type RouterConfig struct {
	UserHandler   *user.Handler
	ChatHandler   *chat.Handler
	ChatWsHandler *chat.WsHandler
	// Add future handlers here, e.g.:
	// FriendHandler *friend.Handler
}

func InitRouter(cfg *RouterConfig) {
	r = gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:5500"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.GET("/", func(c *gin.Context) {
		//time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	// This route is always accessible.
	r.GET("/api/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from a public endpoint! You don't need to be authenticated to see this."})
	})

	// This route is only accessible if the user has a valid access_token.
	r.GET("/api/private", middleware.EnsureValidToken(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from a private endpoint! You need to be authenticated to see this."})
	})

	// All routes in this group are protected
	userRoutes := r.Group("/api/users")
	userRoutes.Use(middleware.EnsureValidToken())
	{
		// Define your routes here, e.g.
		userRoutes.GET("/check", cfg.UserHandler.HandleOAuth2Callback)
		userRoutes.POST("/register", cfg.UserHandler.CreateUser)
		// etc...
	}

	chatRoutes := r.Group("/api/chats")
	chatRoutes.Use(middleware.EnsureValidToken())
	{
		// Define your routes here, e.g.
		chatRoutes.POST("/", cfg.ChatHandler.CreateChatRoom)
		chatRoutes.GET("/:id", cfg.ChatHandler.GetChatRoom)
		// etc...
	}

	// WebSocket api endpoints

	chatWsRoutes := r.Group("/ws/chats")
	//	chatWsRoutes.Use(middleware.EnsureValidToken())
	{
		//chatWsRoutes.POST("/", cfg.ChatWsHandler.WsCreateChatRoom)
		chatWsRoutes.GET("/:id", cfg.ChatWsHandler.JoinRoom)
	}

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
	// kill (no param) default send syscanll
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 1 second.

	<-ctx.Done()
	log.Println("timeout of 1 second")

	log.Println("Server exiting")
}
