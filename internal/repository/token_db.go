package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hantdev/chorus-controller/internal/db"
	"github.com/hantdev/chorus-controller/internal/domain"
	"gorm.io/gorm"
)

type TokenDBRepository struct{}

func NewTokenDBRepository() *TokenDBRepository {
	return &TokenDBRepository{}
}

// Token CRUD operations
func (r *TokenDBRepository) Create(ctx context.Context, token *domain.TokenInfo) error {
	return db.DB().WithContext(ctx).Create(token).Error
}

func (r *TokenDBRepository) List(ctx context.Context) ([]domain.TokenInfo, error) {
	var tokens []domain.TokenInfo
	err := db.DB().WithContext(ctx).Order("created_at desc").Find(&tokens).Error
	return tokens, err
}

func (r *TokenDBRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.TokenInfo, error) {
	var token domain.TokenInfo
	err := db.DB().WithContext(ctx).Where("id = ?", id).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *TokenDBRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.TokenInfo, error) {
	var token domain.TokenInfo
	err := db.DB().WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *TokenDBRepository) GetBySystemName(ctx context.Context, name string, token *domain.TokenInfo) error {
	err := db.DB().WithContext(ctx).Where("name = ?", name).First(token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (r *TokenDBRepository) Update(ctx context.Context, token *domain.TokenInfo) error {
	return db.DB().WithContext(ctx).Save(token).Error
}

func (r *TokenDBRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	result := db.DB().WithContext(ctx).Where("id = ?", id).Delete(&domain.TokenInfo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
