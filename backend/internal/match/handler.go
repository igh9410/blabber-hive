package match

import (
	"backend/internal/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) EnqueueUser(c *gin.Context) {
	// Extract UserID from the context
	userID, err := common.UserIDValidator(c.MustGet("user_id"))
	if err != nil {
		log.Printf("Error occured with user ID %v: %v", userID, err.Error())
	}

	res, err := h.Service.EnqueueUser(c.Request.Context(), userID)

	if err != nil {
		log.Printf("Error occurred with user ID %v: %v", userID, err.Error())
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Handler function
func (h *Handler) DequeueUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Convert userID to UUID if necessary, handle errors

	err := h.Service.DequeueUser(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error occurred with dequeuing user ID %v: %v", userID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing user from queue"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "User removed from queue"})
}
