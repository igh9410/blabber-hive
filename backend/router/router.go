package router

import (
	"github.com/igh9410/blabber-hive/backend/api/middleware"
	myPrometheus "github.com/igh9410/blabber-hive/backend/infra/prometheus"

	"net/http"

	"github.com/igh9410/blabber-hive/backend/internal/api"
	"github.com/igh9410/blabber-hive/backend/internal/chat"
	"github.com/igh9410/blabber-hive/backend/internal/match"
	"github.com/igh9410/blabber-hive/backend/internal/server"
	"github.com/igh9410/blabber-hive/backend/internal/service"
	"github.com/igh9410/blabber-hive/backend/internal/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var r *gin.Engine

type RouterConfig struct {
	UserHandler   *user.Handler
	ChatHandler   *chat.Handler
	ChatWsHandler *chat.WsHandler
	MatchHandler  *match.Handler
	// Add future handlers here, e.g.:
	// FriendHandler *friend.Handler
}

func InitRouter(cfg *RouterConfig) *gin.Engine {
	r = gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:5500", "https://blabberhive.com"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	// Allow all origins for swagger UI
	swagger.Servers = nil

	r.StaticFile("/openapi.yaml", "./api/openapi.yaml")

	r.GET("/", func(c *gin.Context) {
		//time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	myPrometheus.InitPrometheusMetrics()
	// Expose the registered metrics via HTTP.
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(myPrometheus.Reg, promhttp.HandlerOpts{})))

	// Add Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.yaml")))

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
		chatRoutes.POST("/:id", cfg.ChatHandler.JoinChatRoom)
		chatRoutes.GET("/", cfg.ChatHandler.GetChatRoomList)
		chatRoutes.GET("/:id/messages", cfg.ChatHandler.GetChatMessages)
		// etc...
	}

	matchRoutes := r.Group("/api/matches")
	matchRoutes.Use(middleware.EnsureValidToken())
	{
		// Define your routes here, e.g.
		matchRoutes.POST("/", cfg.MatchHandler.EnqueueUser)
		matchRoutes.DELETE("/:userId", cfg.MatchHandler.DequeueUser)

		// etc...
	}

	// WebSocket api endpoints

	chatWsRoutes := r.Group("/ws/chats")
	chatWsRoutes.Use(middleware.WebSocketAuthMiddleware())
	{
		chatWsRoutes.GET("/:id", cfg.ChatWsHandler.RegisterClient)
	}

	chatService := service.NewChatService()

	// Create an instance of your handler that implements api.ServerInterface
	handler := api.NewStrictHandler(server.NewAPI(chatService), nil)

	// Register the handlers with Gin
	api.RegisterHandlers(r, handler)

	return r

}
