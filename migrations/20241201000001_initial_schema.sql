-- CreateExtension creates the "uuid-ossp" extension if it does not exist.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CreateTable creates the "storages" table.
CREATE TABLE "storages" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
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
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "storages_pkey" PRIMARY KEY ("id")
);

-- CreateIndex creates the "idx_storages_name" index on the "storages" table.
CREATE UNIQUE INDEX "idx_storages_name" ON "storages"("name");

-- CreateTable creates the "storage_credentials" table.
CREATE TABLE "storage_credentials" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "storage_id" UUID NOT NULL,
    "user" VARCHAR(255) NOT NULL,
    "access_key_id" VARCHAR(255) NOT NULL,
    "secret_access_key" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "storage_credentials_pkey" PRIMARY KEY ("id")
);

-- CreateIndex creates the "idx_storage_credentials_storage_id" index on the "storage_credentials" table.
CREATE INDEX "idx_storage_credentials_storage_id" ON "storage_credentials"("storage_id");

-- CreateIndex creates the "idx_storage_credentials_user_storage" index on the "storage_credentials" table.
CREATE UNIQUE INDEX "idx_storage_credentials_user_storage" ON "storage_credentials"("user", "storage_id");

-- CreateTable creates the "replicate_jobs" table.
CREATE TABLE "replicate_jobs" (
    "id" UUID NOT NULL DEFAULT uuid_generate_v4(),
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

-- AddForeignKey adds a foreign key constraint to the "storage_credentials" table.
ALTER TABLE "storage_credentials" ADD CONSTRAINT "storage_credentials_storage_id_fkey" FOREIGN KEY ("storage_id") REFERENCES "storages"("id") ON DELETE CASCADE ON UPDATE CASCADE;
