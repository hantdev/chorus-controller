-- Create "replicate_job" table
CREATE TABLE "replicate_job" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "user" character varying(255) NOT NULL,
  "bucket" character varying(255) NOT NULL,
  "from" character varying(255) NOT NULL,
  "to" character varying(255) NOT NULL,
  "to_bucket" character varying(255) NOT NULL,
  "status" character varying(64) NOT NULL DEFAULT 'pending',
  PRIMARY KEY ("id")
);
-- Create index "idx_replicate_job_bucket" to table: "replicate_job"
CREATE INDEX "idx_replicate_job_bucket" ON "replicate_job" ("bucket");
-- Create index "idx_replicate_job_from" to table: "replicate_job"
CREATE INDEX "idx_replicate_job_from" ON "replicate_job" ("from");
-- Create index "idx_replicate_job_to" to table: "replicate_job"
CREATE INDEX "idx_replicate_job_to" ON "replicate_job" ("to");
-- Create index "idx_replicate_job_user" to table: "replicate_job"
CREATE INDEX "idx_replicate_job_user" ON "replicate_job" ("user");
-- Create "storage" table
CREATE TABLE "storage" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "name" character varying(255) NOT NULL,
  "address" character varying(1024) NOT NULL,
  "provider" character varying(64) NOT NULL,
  "is_main" boolean NOT NULL,
  "is_secure" boolean NOT NULL,
  "default_region" character varying(128) NOT NULL,
  "health_check_interval_ms" bigint NOT NULL,
  "http_timeout_ms" bigint NOT NULL,
  "rate_limit_enabled" boolean NOT NULL,
  "rate_limit_rpm" integer NOT NULL,
  "user" character varying(255) NOT NULL,
  "access_key_id" character varying(255) NOT NULL,
  "secret_access_key" character varying(255) NOT NULL,
  "description" character varying(500) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_storage_name" to table: "storage"
CREATE UNIQUE INDEX "idx_storage_name" ON "storage" ("name");
-- Create "token_info" table
CREATE TABLE "token_info" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "name" character varying(255) NOT NULL,
  "description" character varying(500) NOT NULL,
  "token_hash" character varying(255) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_system" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_token_info_expires_at" to table: "token_info"
CREATE INDEX "idx_token_info_expires_at" ON "token_info" ("expires_at");
-- Create index "idx_token_info_token_hash" to table: "token_info"
CREATE UNIQUE INDEX "idx_token_info_token_hash" ON "token_info" ("token_hash");
-- Drop "replicate_jobs" table
DROP TABLE "replicate_jobs";
-- Drop "storages" table
DROP TABLE "storages";
-- Drop "token_infos" table
DROP TABLE "token_infos";
