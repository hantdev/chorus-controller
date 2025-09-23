package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/errors"
)

// ErrorHandler handles API errors and returns appropriate HTTP responses
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error: " + err,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		c.Abort()
	})
}

// HandleError handles API errors and returns appropriate HTTP responses
func HandleError(c *gin.Context, err error) {
	if apiErr, ok := err.(*errors.APIError); ok {
		c.JSON(apiErr.Code, gin.H{
			"error": apiErr.Message,
		})
		return
	}

	// Default to internal server error for unknown errors
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
	})
}

// ErrorResponse creates a standardized error response
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
