package user

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

func (h *Handler) HandleOAuth2Callback(c *gin.Context) {

	userEmail, err := common.EmailValidator(c.MustGet("email"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isRegistered, err := h.Service.IsUserRegistered(c.Request.Context(), userEmail)

	if err != nil {
		log.Println("Error checking user registration:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !isRegistered {
		log.Printf("User with email %v is not registered, needs to complete registration process...", userEmail)
		c.JSON(http.StatusConflict, gin.H{"error": "User is not registered."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User is registered."})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := common.UserIDValidator(c)
	if err != nil {
		log.Printf("Error occured with user ID %v: %v", userID, err.Error())
		return
	}

	userEmail, err := common.EmailValidator(c)
	if err != nil {
		log.Printf("Error occured with user email %v: %v", userEmail, err.Error())
		return
	}

	res, err := h.Service.CreateUser(c.Request.Context(), &req, userID, userEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The user email or ID is already registerd.  " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
