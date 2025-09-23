package domain

import (
	"context"

	pb "github.com/clyso/chorus/proto/gen/go/chorus"
	"google.golang.org/protobuf/types/known/emptypb"
)

// WorkerClient defines the interface for worker communication
type WorkerClient interface {
	GetStorages(ctx context.Context) (*pb.GetStoragesResponse, error)
	ListBucketsForReplication(ctx context.Context, req *pb.ListBucketsForReplicationRequest) (*pb.ListBucketsForReplicationResponse, error)
	AddReplication(ctx context.Context, req *pb.AddReplicationRequest) (*emptypb.Empty, error)
	ListReplications(ctx context.Context) (*pb.ListReplicationsResponse, error)
	PauseReplication(ctx context.Context, req *pb.ReplicationRequest) (*emptypb.Empty, error)
	ResumeReplication(ctx context.Context, req *pb.ReplicationRequest) (*emptypb.Empty, error)
	DeleteReplication(ctx context.Context, req *pb.ReplicationRequest) (*emptypb.Empty, error)
	SwitchBucketZeroDowntime(ctx context.Context, req *pb.SwitchBucketZeroDowntimeRequest) (*emptypb.Empty, error)
}

// ReplicationService defines the interface for replication business logic
type ReplicationService interface {
	CreateReplication(ctx context.Context, req *CreateReplicationRequest) error
	ListReplications(ctx context.Context) ([]*pb.Replication, error)
	PauseReplication(ctx context.Context, id *ReplicationIdentifier) error
	ResumeReplication(ctx context.Context, id *ReplicationIdentifier) error
	DeleteReplication(ctx context.Context, id *ReplicationIdentifier) error
	SwitchZeroDowntime(ctx context.Context, id *ReplicationIdentifier) error
}

// StorageService defines the interface for storage business logic
type StorageService interface {
	ListStorages(ctx context.Context) (*pb.GetStoragesResponse, error)
	ListBuckets(ctx context.Context, req *ListBucketsRequest) (*pb.ListBucketsForReplicationResponse, error)
	CreateStorage(ctx context.Context, storage *Storage) error
	CreateStorageFromRequest(ctx context.Context, req *CreateStorageRequest) error
	ListStorageFromDB(ctx context.Context) ([]Storage, error)
	GetStorageByID(ctx context.Context, id string) (*Storage, error)
	UpdateStorageByID(ctx context.Context, id string, req *CreateStorageRequest) error
	DeleteStorageByID(ctx context.Context, id string) error
}

// TokenService defines the interface for token management
type TokenService interface {
	GenerateToken(ctx context.Context, req *TokenRequest) (*TokenResponse, error)
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)
	ValidateSystemToken(ctx context.Context, token string) error
	RevokeToken(ctx context.Context, token string) error
	RevokeTokenByID(ctx context.Context, id string) error
	ListTokens(ctx context.Context) ([]TokenInfoWithValue, error)
	ListTokensWithValues(ctx context.Context) ([]TokenInfoWithValue, error)
	DeleteToken(ctx context.Context, id string) error
}
