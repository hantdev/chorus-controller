package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	WorkerGRPCAddr string
	HTTPPort       int
	PostgresDSN    string
	JWTSecret      string
	JWTExpiry      time.Duration
	Environment    string
	EncryptionKey  string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func New() (*Config, error) {
	// Load .env file if it exists
	if err := loadEnvFile(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	cfg := &Config{
		WorkerGRPCAddr: getenv("WORKER_GRPC_ADDR", "localhost:9670"),
		PostgresDSN:    getenv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable"),
		JWTSecret:      getenv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		Environment:    getenv("ENV", "development"),
		EncryptionKey:  getenv("ENCRYPTION_KEY", "your-encryption-key-change-this-in-production"),
	}

	// Parse HTTP port
	port := getenv("HTTP_PORT", "8081")
	var p int
	_, err := fmt.Sscanf(port, "%d", &p)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}
	cfg.HTTPPort = p

	// Parse JWT expiry
	jwtExpiryStr := getenv("JWT_EXPIRY", "24h")
	jwtExpiry, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY: %w", err)
	}
	cfg.JWTExpiry = jwtExpiry

	return cfg, nil
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() error {
	// Try to load .env file
	if err := godotenv.Load(); err != nil {
		// If .env doesn't exist, try .env.local
		if err := godotenv.Load(".env.local"); err != nil {
			return err
		}
	}
	return nil
}
