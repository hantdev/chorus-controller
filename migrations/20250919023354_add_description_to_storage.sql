-- Modify "replicate_jobs" table
ALTER TABLE "replicate_jobs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid();
-- Modify "storages" table
ALTER TABLE "storages" ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ADD COLUMN "user" character varying(255) NOT NULL, ADD COLUMN "access_key_id" character varying(255) NOT NULL, ADD COLUMN "secret_access_key" character varying(255) NOT NULL, ADD COLUMN "description" character varying(500) NULL;
-- Drop "storage_credentials" table
DROP TABLE "storage_credentials";
