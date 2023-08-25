package chat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateChatRoom(c *gin.Context) {
	var req CreateChatRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error occured")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.Service.CreateChatRoom(c.Request.Context(), &req)
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
	}
	c.JSON(http.StatusOK, res)
}
