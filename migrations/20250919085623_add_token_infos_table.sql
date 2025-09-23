-- Create "token_infos" table
CREATE TABLE "token_infos" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "description" character varying(500) NULL,
  "token_hash" character varying(255) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create index "idx_token_infos_token_hash" to table: "token_infos"
CREATE UNIQUE INDEX "idx_token_infos_token_hash" ON "token_infos" ("token_hash");
