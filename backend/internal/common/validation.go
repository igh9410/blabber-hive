package common

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ChatRoomIDValidator(c *gin.Context) (uuid.UUID, error) {
	// Extract the id parameter from the request context
	chatRoomIDStr := c.Param("id")

	// Validate the UUID format
	chatRoomID, err := uuid.Parse(chatRoomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return uuid.Nil, err
	}
	log.Println("Chat Room ID: ", chatRoomID)
	return chatRoomID, nil
}

func EmailValidator(email interface{}) (string, error) {
	emailStr, ok := email.(string)
	if !ok {
		return "", fmt.Errorf("email is not a valid string")
	}
	if emailStr == "" {
		return "", fmt.Errorf("email is empty")
	}
	// Further email format validation can be added here if necessary.
	return emailStr, nil
}

func UserIDValidator(userID interface{}) (uuid.UUID, error) {
	userIDStr, ok := userID.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id is not a valid string")
	}
	parsedUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user_id is not a valid UUID: %w", err)
	}
	return parsedUUID, nil
}
