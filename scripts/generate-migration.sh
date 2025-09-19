#!/bin/bash

# Script to generate Atlas migration from model changes
# Usage: ./scripts/generate-migration.sh <migration_name>

if [ $# -eq 0 ]; then
    echo "Usage: $0 <migration_name>"
    echo "Example: $0 add_user_table"
    exit 1
fi

MIGRATION_NAME=$1
TIMESTAMP=$(date +"%Y%m%d%H%M%S")
MIGRATION_FILE="${TIMESTAMP}_${MIGRATION_NAME}.sql"

echo "Generating migration: $MIGRATION_NAME"
echo "Migration file: migrations/$MIGRATION_FILE"

# Generate migration using Atlas
GOWORK=off atlas migrate diff "$MIGRATION_NAME" \
    --dir "file://migrations" \
    --to "ent://internal/domain" \
    --dev-url "docker://postgres/16/dev?search_path=public"

if [ $? -eq 0 ]; then
    echo "✅ Migration generated successfully!"
    echo "📁 Check migrations/ directory for the new file"
    echo "🚀 Run 'atlas migrate apply --env local' to apply the migration"
else
    echo "❌ Failed to generate migration"
    exit 1
fi
