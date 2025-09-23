package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/domain"
	"github.com/hantdev/chorus-controller/internal/middleware"
)

// AuthHandler handles authentication-related endpoints
type AuthHandler struct {
	tokenService domain.TokenService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(tokenService domain.TokenService) *AuthHandler {
	return &AuthHandler{
		tokenService: tokenService,
	}
}

// GenerateToken
// @Summary		Generate a new API token
// @Description	Creates a new API token for accessing protected endpoints
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request		body		domain.TokenRequest	true	"Token generation request"
// @Success		201			{object}	domain.TokenResponse
// @Failure		400			{object}	map[string]interface{}
// @Failure		500			{object}	map[string]interface{}
// @Router			/auth/token [post]
func (h *AuthHandler) GenerateToken(c *gin.Context) {
	var req domain.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse(err))
		return
	}

	resp, err := h.tokenService.GenerateToken(c.Request.Context(), &req)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListTokens
// @Summary		List system tokens
// @Description	Returns a list of system tokens with JWT values (no authentication required)
// @Tags			auth
// @Produce		json
// @Success		200	{array}	domain.TokenInfoWithValue
// @Failure		500	{object}	map[string]interface{}
// @Router			/auth/tokens [get]
func (h *AuthHandler) ListTokens(c *gin.Context) {
	tokens, err := h.tokenService.ListTokens(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// ListTokensWithValues
// @Summary		List all tokens with values
// @Description	Returns detailed information about all tokens including JWT values (requires system token)
// @Tags			auth
// @Produce		json
// @Security		TokenAuth
// @Success		200	{array}	domain.TokenInfoWithValue
// @Failure		401	{object}	map[string]interface{}
// @Failure		500	{object}	map[string]interface{}
// @Router			/auth/tokens/detailed [get]
func (h *AuthHandler) ListTokensWithValues(c *gin.Context) {
	tokens, err := h.tokenService.ListTokensWithValues(c.Request.Context())
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RevokeToken
// @Summary		Revoke an API token by ID
// @Description	Disables an API token by its ID
// @Tags			auth
// @Produce		json
// @Security		TokenAuth
// @Param			token_id	query	string	true	"Token ID to revoke"
// @Success		200		{string}	string	"Token revoked successfully"
// @Failure		400		{object}	map[string]interface{}
// @Failure		401		{object}	map[string]interface{}
// @Failure		404		{object}	map[string]interface{}
// @Router			/auth/revoke [post]
func (h *AuthHandler) RevokeToken(c *gin.Context) {
	tokenID := c.Query("token_id")
	if tokenID == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse(errors.New("token_id parameter is required")))
		return
	}

	if err := h.tokenService.RevokeTokenByID(c.Request.Context(), tokenID); err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token revoked successfully"})
}

// DeleteToken
// @Summary		Delete an API token
// @Description	Permanently deletes an API token
// @Tags			auth
// @Produce		json
// @Security		TokenAuth
// @Param			id	path		string	true	"Token ID"
// @Success		200	{string}	string	"Token deleted successfully"
// @Failure		400	{object}	map[string]interface{}
// @Failure		404	{object}	map[string]interface{}
// @Router			/auth/tokens/{id} [delete]
func (h *AuthHandler) DeleteToken(c *gin.Context) {
	id := c.Param("id")
	if err := h.tokenService.DeleteToken(c.Request.Context(), id); err != nil {
		middleware.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token deleted successfully"})
}
