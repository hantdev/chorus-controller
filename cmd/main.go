// Package main Chorus Controller API
//
// Chorus Controller is a control plane for Chorus Worker that manages S3 replication jobs.
//
//	@title			Chorus Controller API
//	@version		1.0
//	@description		A control plane for Chorus Worker that manages S3 replication jobs
//
//	@host		localhost:8081
//	@BasePath	/
//
//	@schemes	http
package main

import (
	"log"

	"github.com/hantdev/chorus-controller/internal/config"
	"github.com/hantdev/chorus-controller/internal/db"
	"github.com/hantdev/chorus-controller/internal/handler"
	"github.com/hantdev/chorus-controller/internal/repository"
	"github.com/hantdev/chorus-controller/internal/server"
	"github.com/hantdev/chorus-controller/internal/service"

	docs "github.com/hantdev/chorus-controller/docs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize DB
	_, err = db.Open(cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	// Note: AutoMigrate is disabled in favor of Atlas migrations
	// Run 'make migrate-up' to apply migrations
	// if err := database.AutoMigrate(&domain.Storage{}, &domain.StorageCredential{}, &domain.ReplicateJob{}); err != nil {
	// 	log.Fatal(err)
	// }

	// Initialize repository layer
	workerRepo := repository.NewWorkerRepository(cfg.WorkerGRPCAddr)

	// Initialize service layer
	replicationService := service.NewReplicationService(workerRepo)
	storageService := service.NewStorageService(workerRepo)

	// Initialize handler layer
	healthHandler := handler.NewHealthHandler()
	storageHandler := handler.NewStorageHandler(storageService)
	replicationHandler := handler.NewReplicationHandler(replicationService)

	// Initialize server
	srv := server.New(healthHandler, storageHandler, replicationHandler, cfg.HTTPPort)

	// Configure Swagger runtime options
	docs.SwaggerInfo.BasePath = "/"

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
