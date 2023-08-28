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
	/*	var req CreateChatRoomReq
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("Error occured")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} */
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

	userEmail, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email does not exist in the context"})
		return
	}

	emailStr := userEmail.(string)

	log.Printf("About to call FindUserByEmail with email: %s", emailStr)
	user, err := h.UserService.FindUserByEmail(c, emailStr)

	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User with email %s not found", emailStr)})
		return
	}
	res, err := h.Service.JoinChatRoomByID(c, chatRoomID, user.ID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error with joining the chatroom with ID %s and user with ID %s", chatRoomID, user.ID)})
		return
	}
	c.JSON(http.StatusOK, res)
	//c.JSON(http.StatusOK, gin.H{"status": "Successfully joined the chat room", "user_id": user.ID})

}
