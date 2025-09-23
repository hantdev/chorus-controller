-- Modify "replicate_jobs" table
ALTER TABLE "replicate_jobs" ALTER COLUMN "id" SET DEFAULT gen_random_uuid();

-- Modify "storages" table - Step 1: Add columns as nullable first
ALTER TABLE "storages" ALTER COLUMN "id" SET DEFAULT gen_random_uuid();
ALTER TABLE "storages" ADD COLUMN "user" character varying(255) NULL;
ALTER TABLE "storages" ADD COLUMN "access_key_id" character varying(255) NULL;
ALTER TABLE "storages" ADD COLUMN "secret_access_key" character varying(255) NULL;
ALTER TABLE "storages" ADD COLUMN "description" character varying(500) NULL;

-- Step 2: Update existing data with default values
UPDATE "storages" SET 
    "user" = 'default-user',
    "access_key_id" = 'default-access-key',
    "secret_access_key" = 'default-secret-key'
WHERE "user" IS NULL OR "access_key_id" IS NULL OR "secret_access_key" IS NULL;

-- Step 3: Make columns NOT NULL
ALTER TABLE "storages" ALTER COLUMN "user" SET NOT NULL;
ALTER TABLE "storages" ALTER COLUMN "access_key_id" SET NOT NULL;
ALTER TABLE "storages" ALTER COLUMN "secret_access_key" SET NOT NULL;

-- Drop "storage_credentials" table
DROP TABLE "storage_credentials";
