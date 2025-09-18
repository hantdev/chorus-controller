package repository

import (
	"context"

	pb "github.com/clyso/chorus/proto/gen/go/chorus"
	"github.com/hantdev/chorus-controller/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// WorkerRepository implements domain.WorkerClient interface
type WorkerRepository struct {
	addr string
}

// NewWorkerRepository creates a new worker repository
func NewWorkerRepository(addr string) *WorkerRepository {
	return &WorkerRepository{addr: addr}
}

// dial creates a gRPC connection to the worker service
func (r *WorkerRepository) dial(ctx context.Context) (pb.ChorusClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(r.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return pb.NewChorusClient(conn), conn, nil
}

// GetStorages retrieves all configured storages
func (r *WorkerRepository) GetStorages(ctx context.Context) (*pb.GetStoragesResponse, error) {
	client, conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.GetStorages(ctx, &emptypb.Empty{})
}

// ListBucketsForReplication retrieves buckets available for replication
func (r *WorkerRepository) ListBucketsForReplication(ctx context.Context, req *pb.ListBucketsForReplicationRequest) (*pb.ListBucketsForReplicationResponse, error) {
	client, conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.ListBucketsForReplication(ctx, req)
}

// AddReplication creates a new replication job
func (r *WorkerRepository) AddReplication(ctx context.Context, req *pb.AddReplicationRequest) (*emptypb.Empty, error) {
	client, conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.AddReplication(ctx, req)
}

// ListReplications retrieves all replication jobs
func (r *WorkerRepository) ListReplications(ctx context.Context) (*pb.ListReplicationsResponse, error) {
	client, conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.ListReplications(ctx, &emptypb.Empty{})
}

// PauseReplication pauses a replication job
func (r *WorkerRepository) PauseReplication(ctx context.Context, req *pb.ReplicationRequest) (*emptypb.Empty, error) {
	client, conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.PauseReplication(ctx, req)
}

// ResumeReplication resumes a replication job
func (r *WorkerRepository) ResumeReplication(ctx context.Context, req *pb.ReplicationRequest) (*emptypb.Empty, error) {
	client, conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.ResumeReplication(ctx, req)
}

// DeleteReplication deletes a replication job
func (r *WorkerRepository) DeleteReplication(ctx context.Context, req *pb.ReplicationRequest) (*emptypb.Empty, error) {
	client, conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.DeleteReplication(ctx, req)
}

// SwitchBucketZeroDowntime switches buckets without downtime
func (r *WorkerRepository) SwitchBucketZeroDowntime(ctx context.Context, req *pb.SwitchBucketZeroDowntimeRequest) (*emptypb.Empty, error) {
	client, conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.SwitchBucketZeroDowntime(ctx, req)
}

// Ensure WorkerRepository implements domain.WorkerClient interface
var _ domain.WorkerClient = (*WorkerRepository)(nil)
