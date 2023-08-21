package user

import (
	"database/sql"
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

func (h *Handler) HandleOAuth2Callback(c *gin.Context) {
	userEmail, exists := c.Get("email")

	if !exists {
		log.Println("Email does not exist")

	}
	emailStr, ok := userEmail.(string)
	if !ok {
		log.Println("Failed to assert userEmail as string")
		return
	}

	isRegistered, err := h.Service.IsUserRegistered(c.Request.Context(), emailStr)
	if err != nil {
		if err == sql.ErrNoRows { // assuming your service returns this error when a user is not found
			c.Redirect(302, "/register")
			return
		}
		log.Println("Error checking user registration:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !isRegistered {
		c.Redirect(302, "/register")
		return
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail, exists := c.Get("email")

	if !exists {
		log.Println("Email does not exist")

	}

	emailStr, ok := userEmail.(string)
	if !ok {
		log.Println("Failed to assert userEmail as string")
		return
	}

	res, err := h.Service.CreateUser(c.Request.Context(), &req, emailStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
