package httpserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hantdev/chorus-controller/internal/http/api"
	"github.com/hantdev/chorus-controller/internal/worker"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	worker *worker.Client
	port   int
}

func New(workerClient *worker.Client, port int) *Server {
	return &Server{worker: workerClient, port: port}
}

func (s *Server) Run() error {
	r := gin.Default()

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	h := api.New(s.worker)
	r.GET("/health", h.Health)

	r.GET("/storages", h.ListStorages)
	r.GET("/buckets", h.ListBuckets)

	// Storage management
	r.POST("/storages", h.CreateStorage)
	r.PATCH("/storages/:name", h.UpdateStorage)
	r.DELETE("/storages/:name", h.DeleteStorage)

	// Replication management
	r.POST("/replications", h.CreateReplication)
	r.GET("/replications", h.ListReplications)
	r.POST("/replications/pause", h.PauseReplication)
	r.POST("/replications/resume", h.ResumeReplication)
	r.DELETE("/replications", h.DeleteReplication)
	r.POST("/replications/switch/zero-downtime", h.SwitchZeroDowntime)

	return r.Run(fmt.Sprintf(":%d", s.port))
}
