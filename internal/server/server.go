package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/domain"
	"github.com/hantdev/chorus-controller/internal/handler"
	"github.com/hantdev/chorus-controller/internal/middleware"
	"github.com/hantdev/chorus-controller/internal/service"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	healthHandler      *handler.HealthHandler
	storageHandler     *handler.StorageHandler
	replicationHandler *handler.ReplicationHandler
	authHandler        *handler.AuthHandler
	tokenService       domain.TokenService
	port               int
}

// New creates a new HTTP server
func New(
	healthHandler *handler.HealthHandler,
	storageHandler *handler.StorageHandler,
	replicationHandler *handler.ReplicationHandler,
	authHandler *handler.AuthHandler,
	tokenService domain.TokenService,
	port int,
) *Server {
	return &Server{
		healthHandler:      healthHandler,
		storageHandler:     storageHandler,
		replicationHandler: replicationHandler,
		authHandler:        authHandler,
		tokenService:       tokenService,
		port:               port,
	}
}

// Initialize ensures system token exists
func (s *Server) Initialize() error {
	// Ensure system token exists
	if tokenService, ok := s.tokenService.(*service.TokenService); ok {
		return tokenService.EnsureSystemToken(context.Background())
	}
	return nil
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
	r.Use(middleware.ErrorHandler())

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Public endpoints (no authentication required)
	r.GET("/health", s.healthHandler.Health)

	// Authentication endpoints (no authentication required)
	r.POST("/auth/token", s.authHandler.GenerateToken)
	r.GET("/auth/tokens", s.authHandler.ListTokens)

	// System token protected endpoints
	systemProtected := r.Group("/")
	systemProtected.Use(middleware.SystemTokenAuth(s.tokenService))
	{
		systemProtected.GET("/auth/tokens/detailed", s.authHandler.ListTokensWithValues)
		systemProtected.POST("/auth/revoke", s.authHandler.RevokeToken)
		systemProtected.DELETE("/auth/tokens/:id", s.authHandler.DeleteToken)
	}

	// Read-only endpoints (no authentication required)
	r.GET("/storages", s.storageHandler.ListStorages)
	r.GET("/buckets", s.storageHandler.ListBuckets)
	r.GET("/storages/db", s.storageHandler.ListStoragesDB)
	r.GET("/storages/:id", s.storageHandler.GetStorage)
	r.GET("/replications", s.replicationHandler.ListReplications)

	// Protected endpoints (authentication required for write operations)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(s.tokenService))
	{
		// Token management endpoints

		// Storage write operations
		protected.POST("/storages", s.storageHandler.CreateStorage)
		protected.PUT("/storages/:id", s.storageHandler.UpdateStorage)
		protected.DELETE("/storages/:id", s.storageHandler.DeleteStorage)

		// Replication write operations
		protected.POST("/replications", s.replicationHandler.CreateReplication)
		protected.POST("/replications/pause", s.replicationHandler.PauseReplication)
		protected.POST("/replications/resume", s.replicationHandler.ResumeReplication)
		protected.DELETE("/replications", s.replicationHandler.DeleteReplication)
		protected.POST("/replications/switch/zero-downtime", s.replicationHandler.SwitchZeroDowntime)
	}

	return r.Run(fmt.Sprintf(":%d", s.port))
}
