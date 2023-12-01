package match

import (
	"backend/internal/common"
	"log"
	"log/slog"
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

func (h *Handler) EnqueUser(c *gin.Context) {
	// Extract UserID from the context
	userID, err := common.UserIDValidator(c.MustGet("user_id"))
	if err != nil {
		log.Printf("Error occured with user ID %v: %v", userID, err.Error())
	}

	res, err := h.Service.EnqueUser(c.Request.Context(), userID)

	if err != nil {
		slog.Error("Error occured with enqueuing user ID %v: %v", userID, err.Error())
	}

	c.JSON(http.StatusCreated, res)
}
