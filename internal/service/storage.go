package service

import (
	"context"
	"time"

	pb "github.com/clyso/chorus/proto/gen/go/chorus"
	"github.com/google/uuid"
	"github.com/hantdev/chorus-controller/internal/domain"
	"github.com/hantdev/chorus-controller/internal/errors"
	"github.com/hantdev/chorus-controller/internal/repository"
	"gorm.io/gorm"
)

// StorageService implements domain.StorageService interface
type StorageService struct {
	workerClient domain.WorkerClient
	storageRepo  *repository.StorageDBRepository
}

// NewStorageService creates a new storage service
func NewStorageService(workerClient domain.WorkerClient) *StorageService {
	return &StorageService{
		workerClient: workerClient,
		storageRepo:  repository.NewStorageDBRepository(),
	}
}

// ListStorages retrieves all configured storages
func (s *StorageService) ListStorages(ctx context.Context) (*pb.GetStoragesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := s.workerClient.GetStorages(ctx)
	if err != nil {
		return nil, errors.NewBadGatewayError("failed to list storages", err)
	}

	return resp, nil
}

// ListBuckets retrieves buckets available for replication
func (s *StorageService) ListBuckets(ctx context.Context, req *domain.ListBucketsRequest) (*pb.ListBucketsForReplicationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	listReq := &pb.ListBucketsForReplicationRequest{
		User:           req.User,
		From:           req.From,
		To:             req.To,
		ShowReplicated: req.ShowReplicated,
	}

	resp, err := s.workerClient.ListBucketsForReplication(ctx, listReq)
	if err != nil {
		return nil, errors.NewBadGatewayError("failed to list buckets", err)
	}

	return resp, nil
}

// CreateStorage persists a storage config
func (s *StorageService) CreateStorage(ctx context.Context, storage *domain.Storage) error {
	return s.storageRepo.Create(ctx, storage)
}

// CreateStorageFromRequest creates a storage from simplified request
func (s *StorageService) CreateStorageFromRequest(ctx context.Context, req *domain.CreateStorageRequest) error {
	// Create storage with default values
	storage := &domain.Storage{
		Name:                  req.Name,
		Address:               req.Address,
		Provider:              req.Provider,
		IsMain:                false,  // Default to false
		IsSecure:              false,  // Default to false
		DefaultRegion:         "",     // Default region
		HealthCheckIntervalMs: 5000,   // Default 5 seconds
		HttpTimeoutMs:         300000, // Default 5 minutes
		RateLimitEnabled:      false,  // Default disabled
		RateLimitRPM:          0,      // Default no limit
		User:                  req.User,
		AccessKeyID:           req.AccessKey,
		SecretAccessKey:       req.SecretKey,
	}

	return s.storageRepo.Create(ctx, storage)
}

// ListStorageFromDB lists storages from DB
func (s *StorageService) ListStorageFromDB(ctx context.Context) ([]domain.Storage, error) {
	return s.storageRepo.List(ctx)
}

// GetStorageByID retrieves a storage by ID
func (s *StorageService) GetStorageByID(ctx context.Context, id string) (*domain.Storage, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid storage ID format", err)
	}

	storage, err := s.storageRepo.GetByID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("storage not found", err)
		}
		return nil, err
	}
	return storage, nil
}

// UpdateStorageByID updates a storage configuration by ID
func (s *StorageService) UpdateStorageByID(ctx context.Context, id string, req *domain.CreateStorageRequest) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid storage ID format", err)
	}

	// First, get the existing storage
	existingStorage, err := s.storageRepo.GetByID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("storage not found", err)
		}
		return err
	}

	// Update the storage with new values
	existingStorage.Name = req.Name
	existingStorage.Address = req.Address
	existingStorage.Provider = req.Provider
	existingStorage.DefaultRegion = ""
	existingStorage.HealthCheckIntervalMs = 5000
	existingStorage.HttpTimeoutMs = 300000
	existingStorage.RateLimitEnabled = false
	existingStorage.RateLimitRPM = 0

	// Update credentials
	existingStorage.User = req.User
	existingStorage.AccessKeyID = req.AccessKey
	existingStorage.SecretAccessKey = req.SecretKey

	return s.storageRepo.Update(ctx, existingStorage)
}

// DeleteStorageByID deletes a storage by ID
func (s *StorageService) DeleteStorageByID(ctx context.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid storage ID format", err)
	}

	err = s.storageRepo.DeleteByID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("storage not found", err)
		}
		return err
	}
	return nil
}

// Ensure StorageService implements domain.StorageService interface
var _ domain.StorageService = (*StorageService)(nil)
