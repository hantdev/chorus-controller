package api

import (
	pb "github.com/clyso/chorus/proto/gen/go/chorus"
	"github.com/gin-gonic/gin"
)

func errJSON(err error) gin.H { return gin.H{"error": err.Error()} }

func mapUpsertToPB(in UpsertStorage) *pb.UpsertStorageRequest {
	provider := pb.Storage_Other
	switch in.Provider {
	case "Ceph":
		provider = pb.Storage_Ceph
	case "Minio":
		provider = pb.Storage_Minio
	case "AWS":
		provider = pb.Storage_AWS
	case "GCS":
		provider = pb.Storage_GCS
	case "Alibaba":
		provider = pb.Storage_Alibaba
	case "Cloudflare":
		provider = pb.Storage_Cloudflare
	case "DigitalOcean":
		provider = pb.Storage_DigitalOcean
	}
	creds := make([]*pb.Credential, 0, len(in.Credentials))
	for alias, c := range in.Credentials {
		creds = append(creds, &pb.Credential{Alias: alias, AccessKey: c.AccessKey, SecretKey: c.SecretKey})
	}
	return &pb.UpsertStorageRequest{
		Name:                in.Name,
		IsMain:              in.IsMain,
		Address:             in.Address,
		Provider:            provider,
		Credentials:         creds,
		IsSecure:            in.IsSecure,
		DefaultRegion:       in.DefaultRegion,
		HealthCheckInterval: in.HealthCheckInterval,
		HttpTimeout:         in.HttpTimeout,
		RateLimit:           &pb.RateLimitConfig{Enable: in.RateLimit.Enable, Rpm: in.RateLimit.Rpm},
	}
}
