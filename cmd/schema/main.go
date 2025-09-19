package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Export schema to SQL file
	schemaFile := "schema.sql"
	if len(os.Args) > 1 {
		schemaFile = os.Args[1]
	}

	// Write schema to file
	file, err := os.Create(schemaFile)
	if err != nil {
		log.Fatal("Failed to create schema file:", err)
	}
	defer file.Close()

	// Generate PostgreSQL schema directly from GORM models
	generatePostgresSchema(file)

	fmt.Printf("Schema exported to %s\n", schemaFile)
}

func generatePostgresSchema(file *os.File) {
	// Create storages table
	file.WriteString(`-- CreateExtension creates the "uuid-ossp" extension if it does not exist.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CreateTable creates the "storages" table.
CREATE TABLE "storages" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "address" VARCHAR(1024) NOT NULL,
    "provider" VARCHAR(64) NOT NULL,
    "is_main" BOOLEAN NOT NULL DEFAULT false,
    "is_secure" BOOLEAN NOT NULL DEFAULT false,
    "default_region" VARCHAR(128),
    "health_check_interval_ms" BIGINT NOT NULL DEFAULT 0,
    "http_timeout_ms" BIGINT NOT NULL DEFAULT 0,
    "rate_limit_enabled" BOOLEAN NOT NULL DEFAULT false,
    "rate_limit_rpm" INTEGER NOT NULL DEFAULT 0,
    "user" VARCHAR(255) NOT NULL,
    "access_key_id" VARCHAR(255) NOT NULL,
    "secret_access_key" VARCHAR(255) NOT NULL,
    "description" VARCHAR(500),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "storages_pkey" PRIMARY KEY ("id")
);

-- CreateIndex creates the "idx_storages_name" index on the "storages" table.
CREATE UNIQUE INDEX "idx_storages_name" ON "storages"("name");

`)

	// Create replicate_jobs table
	file.WriteString(`-- CreateTable creates the "replicate_jobs" table.
CREATE TABLE "replicate_jobs" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user" VARCHAR(255) NOT NULL,
    "bucket" VARCHAR(255) NOT NULL,
    "from" VARCHAR(255) NOT NULL,
    "to" VARCHAR(255) NOT NULL,
    "to_bucket" VARCHAR(255),
    "status" VARCHAR(64) NOT NULL DEFAULT 'pending',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "replicate_jobs_pkey" PRIMARY KEY ("id")
);

-- CreateIndex creates the "idx_replicate_jobs_user" index on the "replicate_jobs" table.
CREATE INDEX "idx_replicate_jobs_user" ON "replicate_jobs"("user");

-- CreateIndex creates the "idx_replicate_jobs_bucket" index on the "replicate_jobs" table.
CREATE INDEX "idx_replicate_jobs_bucket" ON "replicate_jobs"("bucket");

-- CreateIndex creates the "idx_replicate_jobs_from" index on the "replicate_jobs" table.
CREATE INDEX "idx_replicate_jobs_from" ON "replicate_jobs"("from");

-- CreateIndex creates the "idx_replicate_jobs_to" index on the "replicate_jobs" table.
CREATE INDEX "idx_replicate_jobs_to" ON "replicate_jobs"("to");

`)
}
