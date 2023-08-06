package router

import (
	"backend/api/middleware"
	"backend/internal/user"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler) {
	r = gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// This route is always accessible.
	r.GET("/api/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from a public endpoint! You don't need to be authenticated to see this."})
	})

	// This route is only accessible if the user has a valid access_token.
	r.GET("/api/private", middleware.EnsureValidToken(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from a private endpoint! You need to be authenticated to see this."})
	})

	// All routes in this group are protected
	userRoutes := r.Group("/api/user")
	userRoutes.Use(middleware.EnsureValidToken())
	{
		// Define your routes here, e.g.
		userRoutes.GET("/profile", userHandler.CreateUser)
		userRoutes.POST("/update", userHandler.CreateUser)
		// etc...
	}

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
	log.Print("Server listening on http://localhost:8080")
}

func Start(addr string) error {
	return r.Run(addr)
}
