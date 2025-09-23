package domain

import (
	"time"

	"github.com/google/uuid"
)

// Domain models for the chorus controller

// CreateReplicationRequest represents a request to create a new replication job
type CreateReplicationRequest struct {
	User     string   `json:"user" binding:"required"`
	From     string   `json:"from" binding:"required"`
	To       string   `json:"to" binding:"required"`
	Buckets  []string `json:"buckets"`
	ToBucket string   `json:"to_bucket"`
	AgentURL string   `json:"agent_url"`
}

// CreateStorageRequest represents a simplified request to create a storage
type CreateStorageRequest struct {
	Name      string `json:"name" binding:"required" example:"my-storage"`
	Address   string `json:"address" binding:"required" example:"http://localhost:9000"`
	Provider  string `json:"provider" binding:"required" example:"minio"`
	User      string `json:"user" binding:"required" example:"myuser"`
	AccessKey string `json:"access_key" binding:"required" example:"AKIA123"`
	SecretKey string `json:"secret_key" binding:"required" example:"SECRET123"`
}

// ReplicationIdentifier represents a replication job identifier
type ReplicationIdentifier struct {
	User     string `json:"user" binding:"required"`
	Bucket   string `json:"bucket" binding:"required"`
	From     string `json:"from" binding:"required"`
	To       string `json:"to" binding:"required"`
	ToBucket string `json:"to_bucket"`
}

// ListBucketsRequest represents parameters for listing buckets
type ListBucketsRequest struct {
	User           string `form:"user"`
	From           string `form:"from"`
	To             string `form:"to"`
	ShowReplicated bool   `form:"show_replicated"`
}

// Storage represents a storage configuration persisted in DB
// Mirrors fields from chorus-worker's s3.Storage and adds Name
// Each storage has one user with embedded credentials
type Storage struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name                  string    `gorm:"uniqueIndex;size:255;not null" json:"name"`
	Address               string    `gorm:"size:1024;not null" json:"address"`
	Provider              string    `gorm:"size:64;not null" json:"provider"`
	IsMain                bool      `json:"is_main"`
	IsSecure              bool      `json:"is_secure"`
	DefaultRegion         string    `gorm:"size:128" json:"default_region"`
	HealthCheckIntervalMs int64     `json:"health_check_interval_ms"`
	HttpTimeoutMs         int64     `json:"http_timeout_ms"`
	RateLimitEnabled      bool      `json:"rate_limit_enabled"`
	RateLimitRPM          int       `json:"rate_limit_rpm"`
	User                  string    `gorm:"size:255;not null" json:"user"`
	AccessKeyID           string    `gorm:"size:255;not null" json:"access_key_id"`
	SecretAccessKey       string    `gorm:"size:255;not null" json:"secret_access_key"`
	Description           string    `gorm:"size:500" json:"description"`
}

// TableName returns the table name for Storage
func (Storage) TableName() string {
	return "storage"
}

// ReplicateJob represents a replication job persisted in DB
// Each bucket under a user/from/to is a row; ToBucket can be empty for same name.
type ReplicateJob struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	User     string    `gorm:"size:255;index;not null" json:"user"`
	Bucket   string    `gorm:"size:255;index;not null" json:"bucket"`
	From     string    `gorm:"size:255;index;not null" json:"from"`
	To       string    `gorm:"size:255;index;not null" json:"to"`
	ToBucket string    `gorm:"size:255" json:"to_bucket"`
	Status   string    `gorm:"size:64;default:'pending'" json:"status"`
}

// TableName returns the table name for ReplicateJob
func (ReplicateJob) TableName() string {
	return "replicate_job"
}

// TokenRequest represents a request to generate a token
type TokenRequest struct {
	Name        string     `json:"name" binding:"required" example:"api-client"`
	Description string     `json:"description" example:"Token for API access"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

// TokenResponse represents a token generation response
type TokenResponse struct {
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenInfo represents token information stored in database
type TokenInfo struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description string     `gorm:"size:500" json:"description"`
	TokenHash   string     `gorm:"size:255;not null;uniqueIndex" json:"-"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	IsSystem    bool       `gorm:"default:false" json:"is_system"`
	ExpiresAt   *time.Time `gorm:"index" json:"expires_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for TokenInfo
func (TokenInfo) TableName() string {
	return "token_info"
}

// TokenInfoWithValue represents token information with JWT token value for API responses
type TokenInfoWithValue struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Token       string     `json:"token"`
	IsActive    bool       `json:"is_active"`
	IsSystem    bool       `json:"is_system"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
