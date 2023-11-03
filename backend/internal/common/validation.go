package common

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
func EmailValidator(c *gin.Context) (string, error) {
	// Extract the email parameter from the request context
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email does not exist in the context"})
		return "", nil
	}
	log.Println("Email: ", email)
	return email.(string), nil
}*/

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

/*
func UserIDValidator(c *gin.Context) (uuid.UUID, error) {

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID does not exist in the context"})
		return uuid.Nil, nil
	}

	// Validate the UUID format
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return uuid.Nil, err
	}
	log.Println("User ID: ", userID)
	return userID, nil
} */

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
