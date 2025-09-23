package service

import (
	"context"
	"time"

	pb "github.com/clyso/chorus/proto/gen/go/chorus"
	"github.com/google/uuid"
	"github.com/hantdev/chorus-controller/internal/crypto"
	"github.com/hantdev/chorus-controller/internal/domain"
	"github.com/hantdev/chorus-controller/internal/errors"
	"github.com/hantdev/chorus-controller/internal/repository"
	"gorm.io/gorm"
)

// StorageService implements domain.StorageService interface
type StorageService struct {
	workerClient domain.WorkerClient
	storageRepo  *repository.StorageDBRepository
	crypto       *crypto.Crypto
}

// NewStorageService creates a new storage service
func NewStorageService(workerClient domain.WorkerClient, encryptionKey string) *StorageService {
	crypto, err := crypto.New(encryptionKey)
	if err != nil {
		// In production, you might want to handle this error more gracefully
		panic("failed to initialize crypto: " + err.Error())
	}

	return &StorageService{
		workerClient: workerClient,
		storageRepo:  repository.NewStorageDBRepository(),
		crypto:       crypto,
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
	// Encrypt sensitive data before saving
	if err := s.encryptStorage(storage); err != nil {
		return err
	}
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

	// Encrypt sensitive data before saving
	if err := s.encryptStorage(storage); err != nil {
		return err
	}

	return s.storageRepo.Create(ctx, storage)
}

// ListStorageFromDB lists storages from DB
func (s *StorageService) ListStorageFromDB(ctx context.Context) ([]domain.Storage, error) {
	storages, err := s.storageRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive data for each storage
	for i := range storages {
		if err := s.decryptStorage(&storages[i]); err != nil {
			return nil, err
		}
	}

	return storages, nil
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

	// Decrypt sensitive data before returning
	if err := s.decryptStorage(storage); err != nil {
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

	// Encrypt sensitive data before saving
	if err := s.encryptStorage(existingStorage); err != nil {
		return err
	}

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

// encryptStorage encrypts sensitive fields before saving to database
func (s *StorageService) encryptStorage(storage *domain.Storage) error {
	if storage.SecretAccessKey != "" {
		encrypted, err := s.crypto.Encrypt(storage.SecretAccessKey)
		if err != nil {
			return errors.NewInternalServerError("failed to encrypt secret access key", err)
		}
		storage.SecretAccessKey = encrypted
	}
	return nil
}

// decryptStorage decrypts sensitive fields after loading from database
func (s *StorageService) decryptStorage(storage *domain.Storage) error {
	if storage.SecretAccessKey != "" {
		decrypted, err := s.crypto.Decrypt(storage.SecretAccessKey)
		if err != nil {
			return errors.NewInternalServerError("failed to decrypt secret access key", err)
		}
		storage.SecretAccessKey = decrypted
	}
	return nil
}

// Ensure StorageService implements domain.StorageService interface
var _ domain.StorageService = (*StorageService)(nil)
