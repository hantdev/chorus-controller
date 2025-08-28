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
	httpserver "github.com/hantdev/chorus-controller/internal/http"
	"github.com/hantdev/chorus-controller/internal/worker"

	docs "github.com/hantdev/chorus-controller/docs"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	w := worker.New(cfg.WorkerGRPCAddr)
	srv := httpserver.New(w, cfg.HTTPPort)

	// Configure Swagger runtime options
	docs.SwaggerInfo.BasePath = "/"

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
