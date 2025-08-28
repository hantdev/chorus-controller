package config

import (
	"fmt"
	"os"
)

type Config struct {
	WorkerGRPCAddr string
	HTTPPort       int
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func New() (*Config, error) {
	cfg := &Config{
		WorkerGRPCAddr: getenv("WORKER_GRPC_ADDR", "localhost:9670"),
	}
	port := getenv("HTTP_PORT", "8081")
	var p int
	_, err := fmt.Sscanf(port, "%d", &p)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}
	cfg.HTTPPort = p
	return cfg, nil
}
