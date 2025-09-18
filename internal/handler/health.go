package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health
// @Summary		Health check
// @Description	Returns the health status of the controller
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]interface{}
// @Router			/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
