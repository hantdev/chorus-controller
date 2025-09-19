#!/bin/bash

# Script to apply Atlas migrations
# Usage: ./scripts/apply-migration.sh [environment]

ENV=${1:-local}

echo "Applying migrations to environment: $ENV"

# Apply migrations using Atlas
GOWORK=off atlas migrate apply --env "$ENV"

if [ $? -eq 0 ]; then
    echo "✅ Migrations applied successfully!"
else
    echo "❌ Failed to apply migrations"
    exit 1
fi
