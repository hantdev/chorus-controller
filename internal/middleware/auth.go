package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/domain"
)

// AuthMiddleware creates a middleware for JWT token authentication
func AuthMiddleware(tokenService domain.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("Authorization header is required")))
			c.Abort()
			return
		}

		// Check if it's a Token
		if !strings.HasPrefix(authHeader, "Token ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("Invalid authorization header format. Expected 'Token <token>'")))
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Token ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("Token is required")))
			c.Abort()
			return
		}

		// Validate token
		tokenInfo, err := tokenService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse(err))
			c.Abort()
			return
		}

		// Store token info in context for use in handlers
		c.Set("token_info", tokenInfo)
		c.Next()
	}
}

// OptionalAuthMiddleware creates a middleware for optional JWT token authentication
// This allows endpoints to work with or without authentication
func OptionalAuthMiddleware(tokenService domain.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Check if it's a Token
		if !strings.HasPrefix(authHeader, "Token ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("Invalid authorization header format. Expected 'Token <token>'")))
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Token ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("Token is required")))
			c.Abort()
			return
		}

		// Validate token
		tokenInfo, err := tokenService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse(err))
			c.Abort()
			return
		}

		// Store token info in context for use in handlers
		c.Set("token_info", tokenInfo)
		c.Next()
	}
}

// GetTokenInfoFromContext extracts token info from gin context
func GetTokenInfoFromContext(c *gin.Context) (*domain.TokenInfo, error) {
	tokenInfo, exists := c.Get("token_info")
	if !exists {
		return nil, errors.New("token info not found in context")
	}

	token, ok := tokenInfo.(*domain.TokenInfo)
	if !ok {
		return nil, errors.New("invalid token info type")
	}

	return token, nil
}
