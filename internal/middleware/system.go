package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/domain"
)

// SystemTokenAuth middleware validates system token
func SystemTokenAuth(tokenService domain.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("Authorization header is required")))
			c.Abort()
			return
		}

		// Expect format: "Token <token>"
		if !strings.HasPrefix(authHeader, "Token ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("Invalid authorization header format. Expected 'Token <token>'")))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Token ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("Token is required")))
			c.Abort()
			return
		}
		if err := tokenService.ValidateSystemToken(c.Request.Context(), token); err != nil {
			HandleError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
