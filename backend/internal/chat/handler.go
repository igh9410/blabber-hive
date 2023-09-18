package chat

import (
	"backend/internal/user"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service
	UserService user.Service
}

func NewHandler(s Service, u user.Service) *Handler {
	return &Handler{
		Service:     s,
		UserService: u,
	}
}

func (h *Handler) CreateChatRoom(c *gin.Context) {

	res, err := h.Service.CreateChatRoom(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetChatRoom(c *gin.Context) {
	// Extract the id parameter from the request context
	chatRoomIDStr := c.Param("id")

	// Validate the UUID format
	chatRoomID, err := uuid.Parse(chatRoomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	res, err := h.Service.GetChatRoomByID(c.Request.Context(), chatRoomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Chat ID %s not found", chatRoomID)})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) JoinChatRoom(c *gin.Context) {

	// Extract the id parameter from the request context
	chatRoomIDStr := c.Param("id")

	// Validate the UUID format
	chatRoomID, err := uuid.Parse(chatRoomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID does not exist in the context"})
		return
	}

	// Validate the UUID format
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	log.Println("User ID: ", userID)

	res, err := h.Service.JoinChatRoomByID(c, chatRoomID, userID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error with joining the chatroom with ID %s and user with ID %s", chatRoomID, userID)})
		return
	}
	c.JSON(http.StatusOK, res)

}
