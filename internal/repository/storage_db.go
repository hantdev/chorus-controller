package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hantdev/chorus-controller/internal/db"
	"github.com/hantdev/chorus-controller/internal/domain"
	"gorm.io/gorm"
)

type StorageDBRepository struct{}

type ReplicateJobDBRepository struct{}

func NewStorageDBRepository() *StorageDBRepository           { return &StorageDBRepository{} }
func NewReplicateJobDBRepository() *ReplicateJobDBRepository { return &ReplicateJobDBRepository{} }

// Storage CRUD
func (r *StorageDBRepository) Create(ctx context.Context, s *domain.Storage) error {
	return db.DB().WithContext(ctx).Create(s).Error
}

func (r *StorageDBRepository) List(ctx context.Context) ([]domain.Storage, error) {
	var items []domain.Storage
	err := db.DB().WithContext(ctx).Preload("Credentials").Order("name asc").Find(&items).Error
	return items, err
}

func (r *StorageDBRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Storage, error) {
	var s domain.Storage
	err := db.DB().WithContext(ctx).Preload("Credentials").Where("id = ?", id).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *StorageDBRepository) Update(ctx context.Context, s *domain.Storage) error {
	return db.DB().WithContext(ctx).Save(s).Error
}

func (r *StorageDBRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	result := db.DB().WithContext(ctx).Where("id = ?", id).Delete(&domain.Storage{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ReplicateJob CRUD
func (r *ReplicateJobDBRepository) Create(ctx context.Context, j *domain.ReplicateJob) error {
	return db.DB().WithContext(ctx).Create(j).Error
}

func (r *ReplicateJobDBRepository) List(ctx context.Context) ([]domain.ReplicateJob, error) {
	var items []domain.ReplicateJob
	err := db.DB().WithContext(ctx).Order("id desc").Find(&items).Error
	return items, err
}

func (r *ReplicateJobDBRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return db.DB().WithContext(ctx).Where("id = ?", id).Delete(&domain.ReplicateJob{}).Error
}
