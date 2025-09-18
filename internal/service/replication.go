package service

import (
	"context"
	"time"

	pb "github.com/clyso/chorus/proto/gen/go/chorus"
	"github.com/hantdev/chorus-controller/internal/domain"
	"github.com/hantdev/chorus-controller/internal/errors"
	"github.com/hantdev/chorus-controller/internal/repository"
)

// ReplicationService implements domain.ReplicationService interface
type ReplicationService struct {
	workerClient     domain.WorkerClient
	replicateJobRepo *repository.ReplicateJobDBRepository
}

// NewReplicationService creates a new replication service
func NewReplicationService(workerClient domain.WorkerClient) *ReplicationService {
	return &ReplicationService{
		workerClient:     workerClient,
		replicateJobRepo: repository.NewReplicateJobDBRepository(),
	}
}

// CreateReplication creates a new replication job
func (s *ReplicationService) CreateReplication(ctx context.Context, req *domain.CreateReplicationRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	addReq := &pb.AddReplicationRequest{
		User:            req.User,
		From:            req.From,
		To:              req.To,
		Buckets:         req.Buckets,
		IsForAllBuckets: len(req.Buckets) == 0,
		ToBucket:        req.ToBucket,
	}

	if req.AgentURL != "" {
		addReq.AgentUrl = &req.AgentURL
	}

	_, err := s.workerClient.AddReplication(ctx, addReq)
	if err != nil {
		return errors.NewBadGatewayError("failed to create replication", err)
	}

	// Persist ReplicateJob(s) into DB for tracking
	toBucket := req.ToBucket
	if toBucket == "" {
		// when empty, destination bucket equals source bucket
		toBucket = ""
	}
	if len(req.Buckets) > 0 {
		for _, b := range req.Buckets {
			jb := &domain.ReplicateJob{
				User:     req.User,
				Bucket:   b,
				From:     req.From,
				To:       req.To,
				ToBucket: toBucket,
				Status:   "created",
			}
			_ = s.replicateJobRepo.Create(ctx, jb)
		}
	} else {
		// represent all buckets with empty bucket field
		jb := &domain.ReplicateJob{
			User:   req.User,
			Bucket: "",
			From:   req.From,
			To:     req.To,
			// ToBucket left empty to signify same name
			Status: "created",
		}
		_ = s.replicateJobRepo.Create(ctx, jb)
	}

	return nil
}

// ListReplications retrieves all replication jobs
func (s *ReplicationService) ListReplications(ctx context.Context) ([]*pb.Replication, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := s.workerClient.ListReplications(ctx)
	if err != nil {
		return nil, errors.NewBadGatewayError("failed to list replications", err)
	}

	return resp.Replications, nil
}

// PauseReplication pauses a replication job
func (s *ReplicationService) PauseReplication(ctx context.Context, id *domain.ReplicationIdentifier) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := s.buildReplicationRequest(id)
	_, err := s.workerClient.PauseReplication(ctx, req)
	if err != nil {
		return errors.NewBadGatewayError("failed to pause replication", err)
	}

	return nil
}

// ResumeReplication resumes a replication job
func (s *ReplicationService) ResumeReplication(ctx context.Context, id *domain.ReplicationIdentifier) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := s.buildReplicationRequest(id)
	_, err := s.workerClient.ResumeReplication(ctx, req)
	if err != nil {
		return errors.NewBadGatewayError("failed to resume replication", err)
	}

	return nil
}

// DeleteReplication deletes a replication job
func (s *ReplicationService) DeleteReplication(ctx context.Context, id *domain.ReplicationIdentifier) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := s.buildReplicationRequest(id)
	_, err := s.workerClient.DeleteReplication(ctx, req)
	if err != nil {
		return errors.NewBadGatewayError("failed to delete replication", err)
	}

	return nil
}

// SwitchZeroDowntime switches buckets without downtime
func (s *ReplicationService) SwitchZeroDowntime(ctx context.Context, id *domain.ReplicationIdentifier) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := s.buildReplicationRequest(id)
	switchReq := &pb.SwitchBucketZeroDowntimeRequest{
		ReplicationId: req,
	}

	_, err := s.workerClient.SwitchBucketZeroDowntime(ctx, switchReq)
	if err != nil {
		return errors.NewBadGatewayError("failed to switch buckets", err)
	}

	return nil
}

// buildReplicationRequest builds a pb.ReplicationRequest from domain.ReplicationIdentifier
func (s *ReplicationService) buildReplicationRequest(id *domain.ReplicationIdentifier) *pb.ReplicationRequest {
	toBucket := id.ToBucket
	if toBucket == "" {
		toBucket = id.Bucket
	}

	return &pb.ReplicationRequest{
		User:     id.User,
		Bucket:   id.Bucket,
		From:     id.From,
		To:       id.To,
		ToBucket: toBucket,
	}
}

// Ensure ReplicationService implements domain.ReplicationService interface
var _ domain.ReplicationService = (*ReplicationService)(nil)
