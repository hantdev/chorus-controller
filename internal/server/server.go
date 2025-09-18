package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/handler"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	healthHandler      *handler.HealthHandler
	storageHandler     *handler.StorageHandler
	replicationHandler *handler.ReplicationHandler
	port               int
}

// New creates a new HTTP server
func New(
	healthHandler *handler.HealthHandler,
	storageHandler *handler.StorageHandler,
	replicationHandler *handler.ReplicationHandler,
	port int,
) *Server {
	return &Server{
		healthHandler:      healthHandler,
		storageHandler:     storageHandler,
		replicationHandler: replicationHandler,
		port:               port,
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins in development
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add error handling middleware
	r.Use(handler.ErrorHandler())

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Health check
	r.GET("/health", s.healthHandler.Health)

	// Storage endpoints
	r.GET("/storages", s.storageHandler.ListStorages)
	r.GET("/buckets", s.storageHandler.ListBuckets)
	// DB-backed storage endpoints
	r.POST("/storages", s.storageHandler.CreateStorage)
	r.GET("/storages/db", s.storageHandler.ListStoragesDB)
	r.GET("/storages/:id", s.storageHandler.GetStorage)
	r.PUT("/storages/:id", s.storageHandler.UpdateStorage)
	r.DELETE("/storages/:id", s.storageHandler.DeleteStorage)

	// Replication management endpoints
	r.POST("/replications", s.replicationHandler.CreateReplication)
	r.GET("/replications", s.replicationHandler.ListReplications)
	r.POST("/replications/pause", s.replicationHandler.PauseReplication)
	r.POST("/replications/resume", s.replicationHandler.ResumeReplication)
	r.DELETE("/replications", s.replicationHandler.DeleteReplication)
	r.POST("/replications/switch/zero-downtime", s.replicationHandler.SwitchZeroDowntime)

	return r.Run(fmt.Sprintf(":%d", s.port))
}
