package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hantdev/chorus-controller/internal/domain"
	"github.com/hantdev/chorus-controller/internal/errors"
	"github.com/hantdev/chorus-controller/internal/repository"
	"gorm.io/gorm"
)

// TokenService implements domain.TokenService interface
type TokenService struct {
	tokenRepo *repository.TokenDBRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// NewTokenService creates a new token service
func NewTokenService(tokenRepo *repository.TokenDBRepository, jwtSecret string, jwtExpiry time.Duration) *TokenService {
	return &TokenService{
		tokenRepo: tokenRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// GenerateToken creates a new API token
func (s *TokenService) GenerateToken(ctx context.Context, req *domain.TokenRequest) (*domain.TokenResponse, error) {
	// Generate a random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, errors.NewInternalServerError("failed to generate token", err)
	}
	tokenString := hex.EncodeToString(tokenBytes)

	// Hash the token for storage
	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])

	// Set default expiration time to 24 hours if not provided
	// System tokens never expire
	var expiresAt *time.Time
	if req.Name == "system" {
		expiresAt = nil // System token never expires
	} else if req.ExpiresAt != nil {
		expiresAt = req.ExpiresAt
	} else {
		defaultExpiry := time.Now().Add(24 * time.Hour)
		expiresAt = &defaultExpiry
	}

	// Create token info
	tokenInfo := &domain.TokenInfo{
		Name:        req.Name,
		Description: req.Description,
		TokenHash:   tokenHash,
		IsActive:    true,
		ExpiresAt:   expiresAt,
	}

	if err := s.tokenRepo.Create(ctx, tokenInfo); err != nil {
		return nil, err
	}

	// Generate JWT token with token info
	jwtToken, _, err := s.generateJWTToken(tokenInfo, tokenString)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to generate JWT token", err)
	}

	return &domain.TokenResponse{
		Token:     jwtToken,
		Name:      tokenInfo.Name,
		ExpiresAt: *expiresAt,
		CreatedAt: tokenInfo.CreatedAt,
	}, nil
}

// ValidateToken validates a JWT token and returns the token info
func (s *TokenService) ValidateToken(ctx context.Context, tokenString string) (*domain.TokenInfo, error) {
	// Parse JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid token", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenIDStr, ok := claims["token_id"].(string)
		if !ok {
			return nil, errors.NewUnauthorizedError("invalid token claims", nil)
		}

		tokenID, err := uuid.Parse(tokenIDStr)
		if err != nil {
			return nil, errors.NewUnauthorizedError("invalid token ID in token", err)
		}

		// Get token info from database
		tokenInfo, err := s.tokenRepo.GetByID(ctx, tokenID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewUnauthorizedError("token not found", nil)
			}
			return nil, err
		}

		// Check if token is still active
		if !tokenInfo.IsActive {
			return nil, errors.NewUnauthorizedError("token is disabled", nil)
		}

		return tokenInfo, nil
	}

	return nil, errors.NewUnauthorizedError("invalid token", nil)
}

// RevokeToken disables a token
func (s *TokenService) RevokeToken(ctx context.Context, tokenString string) error {
	tokenInfo, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return err
	}

	tokenInfo.IsActive = false
	return s.tokenRepo.Update(ctx, tokenInfo)
}

// RevokeTokenByID disables a token by ID
func (s *TokenService) RevokeTokenByID(ctx context.Context, id string) error {
	tokenUUID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid token ID format", err)
	}

	tokenInfo, err := s.tokenRepo.GetByID(ctx, tokenUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("token not found", err)
		}
		return err
	}

	tokenInfo.IsActive = false
	return s.tokenRepo.Update(ctx, tokenInfo)
}

// ListTokens retrieves only system tokens with JWT values (no auth required)
func (s *TokenService) ListTokens(ctx context.Context) ([]domain.TokenInfoWithValue, error) {
	// Only return system tokens with JWT values
	var systemTokens []domain.TokenInfoWithValue
	tokens, err := s.tokenRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		if token.IsSystem {
			// Generate JWT token for system token
			jwtToken, _, err := s.generateJWTToken(&token, "")
			if err != nil {
				continue
			}

			tokenWithValue := domain.TokenInfoWithValue{
				ID:          token.ID,
				Name:        token.Name,
				Description: token.Description,
				Token:       jwtToken,
				IsActive:    token.IsActive,
				IsSystem:    token.IsSystem,
				ExpiresAt:   token.ExpiresAt, // Will be nil for system tokens
				CreatedAt:   token.CreatedAt,
				UpdatedAt:   token.UpdatedAt,
			}
			systemTokens = append(systemTokens, tokenWithValue)
		}
	}

	return systemTokens, nil
}

// ListTokensWithValues retrieves all tokens with their JWT token values
func (s *TokenService) ListTokensWithValues(ctx context.Context) ([]domain.TokenInfoWithValue, error) {
	tokens, err := s.tokenRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []domain.TokenInfoWithValue
	for _, token := range tokens {
		// Generate JWT token for each token info
		// We need to reconstruct the original token string from the hash
		// Since we can't reverse the hash, we'll generate a new JWT with the same claims
		jwtToken, _, err := s.generateJWTToken(&token, "")
		if err != nil {
			// If we can't generate JWT, skip this token
			continue
		}

		tokenWithValue := domain.TokenInfoWithValue{
			ID:          token.ID,
			Name:        token.Name,
			Description: token.Description,
			Token:       jwtToken,
			IsActive:    token.IsActive,
			IsSystem:    token.IsSystem,
			ExpiresAt:   token.ExpiresAt,
			CreatedAt:   token.CreatedAt,
			UpdatedAt:   token.UpdatedAt,
		}
		result = append(result, tokenWithValue)
	}

	return result, nil
}

// DeleteToken deletes a token by ID
func (s *TokenService) DeleteToken(ctx context.Context, id string) error {
	tokenUUID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid token ID format", err)
	}

	// Check if token exists and is not a system token
	tokenInfo, err := s.tokenRepo.GetByID(ctx, tokenUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("token not found", err)
		}
		return err
	}

	// Prevent deletion of system tokens
	if tokenInfo.IsSystem {
		return errors.NewBadRequestError("cannot delete system token", nil)
	}

	err = s.tokenRepo.DeleteByID(ctx, tokenUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("token not found", err)
		}
		return err
	}

	return nil
}

// generateJWTToken creates a JWT token for the given token info
func (s *TokenService) generateJWTToken(tokenInfo *domain.TokenInfo, originalToken string) (string, time.Time, error) {
	var expiresAt time.Time
	var claims jwt.MapClaims

	if tokenInfo.IsSystem {
		// System tokens never expire - don't set exp claim
		claims = jwt.MapClaims{
			"token_id": tokenInfo.ID.String(),
			"name":     tokenInfo.Name,
			"iat":      time.Now().Unix(),
		}
		expiresAt = time.Time{} // Zero time for never-expiring tokens
	} else if tokenInfo.ExpiresAt != nil {
		expiresAt = *tokenInfo.ExpiresAt
		claims = jwt.MapClaims{
			"token_id": tokenInfo.ID.String(),
			"name":     tokenInfo.Name,
			"exp":      expiresAt.Unix(),
			"iat":      time.Now().Unix(),
		}
	} else {
		expiresAt = time.Now().Add(s.jwtExpiry)
		claims = jwt.MapClaims{
			"token_id": tokenInfo.ID.String(),
			"name":     tokenInfo.Name,
			"exp":      expiresAt.Unix(),
			"iat":      time.Now().Unix(),
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// EnsureSystemToken creates or ensures the system token exists
func (s *TokenService) EnsureSystemToken(ctx context.Context) error {
	// Check if system token already exists
	var existingToken domain.TokenInfo
	err := s.tokenRepo.GetBySystemName(ctx, "system", &existingToken)
	if err == nil {
		// System token exists, check if it needs to be updated
		needsUpdate := false

		if !existingToken.IsSystem {
			existingToken.IsSystem = true
			needsUpdate = true
		}

		// If token has expiration, remove it for system tokens
		if existingToken.ExpiresAt != nil {
			existingToken.ExpiresAt = nil
			needsUpdate = true
		}

		if needsUpdate {
			return s.tokenRepo.Update(ctx, &existingToken)
		}
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Create system token
	systemToken := &domain.TokenRequest{
		Name:        "system",
		Description: "Default system token for internal operations",
		ExpiresAt:   nil, // System token never expires
	}

	// Generate the token
	resp, err := s.GenerateToken(ctx, systemToken)
	if err != nil {
		return err
	}

	// Mark as system token
	tokenInfo, err := s.tokenRepo.GetByTokenHash(ctx, s.hashToken(resp.Token))
	if err != nil {
		return err
	}

	tokenInfo.IsSystem = true
	tokenInfo.ExpiresAt = nil // Ensure system token never expires
	return s.tokenRepo.Update(ctx, tokenInfo)
}

// ValidateSystemToken validates if the provided token is a system token
func (s *TokenService) ValidateSystemToken(ctx context.Context, tokenString string) error {
	tokenInfo, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return err
	}

	if !tokenInfo.IsSystem {
		return errors.NewUnauthorizedError("system token required", nil)
	}

	return nil
}

// hashToken creates a hash of the token for storage
func (s *TokenService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Ensure TokenService implements domain.TokenService interface
var _ domain.TokenService = (*TokenService)(nil)
