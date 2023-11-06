package chat

import (
	"backend/internal/common"
	"backend/internal/user"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	// Extract ChatRoomID from the context
	chatRoomID, err := common.ChatRoomIDValidator(c)

	if err != nil {
		log.Printf("Error occured with chat room ID %v: %v", chatRoomID, err.Error())
	}
	res, err := h.Service.GetChatRoomInfoByID(c.Request.Context(), chatRoomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Chat ID %s not found", chatRoomID)})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) JoinChatRoom(c *gin.Context) {
	// Extract ChatRoomID from the context
	chatRoomID, err := common.ChatRoomIDValidator(c)

	if err != nil {
		log.Printf("Error occured with chat room ID %v: %v", chatRoomID, err.Error())
	}

	// Extract UserID from the context
	userID, err := common.UserIDValidator(c)
	if err != nil {
		log.Printf("Error occured with user ID %v: %v", userID, err.Error())
	}

	res, err := h.Service.JoinChatRoomByID(c, chatRoomID, userID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error with joining the chatroom with ID %s and user with ID %s", chatRoomID, userID)})
		return
	}
	c.JSON(http.StatusOK, res)

}

func (h *Handler) GetChatMessages(c *gin.Context) {
	// Extract ChatRoomID from the context
	chatRoomID, err := common.ChatRoomIDValidator(c)
	if err != nil {
		log.Printf("Error occurred with chat room ID: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat room ID"})
		return
	}

	// Extract the cursor and pageSize from query parameters
	cursorStr := c.DefaultQuery("cursor", "") // Cursor is provided by the client
	pageSize := c.DefaultQuery("pageSize", "50")

	// Parse pageSize to an integer
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	// Parse the cursor string to a time.Time
	var cursorTime *time.Time
	if cursorStr != "" {
		parsedCursor, err := time.Parse(time.RFC3339, cursorStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor"})
			return
		}
		cursorTime = &parsedCursor
	}

	// Call the service function to get paginated messages
	messages, err := h.Service.GetPaginatedMessages(c.Request.Context(), chatRoomID, cursorTime, pageSizeInt)
	if err != nil {
		log.Printf("Error fetching messages: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching messages"})
		return
	}

	// Determine the next cursor
	var nextCursor string
	if len(messages) > 0 {
		nextCursor = messages[len(messages)-1].CreatedAt.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages, "nextCursor": nextCursor})
}
